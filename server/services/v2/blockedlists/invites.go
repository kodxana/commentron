package blockedlists

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/lbryio/commentron/server/auth"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func listInvites(r *http.Request, args *commentapi.SharedBlockedListListInvitesArgs, reply *commentapi.SharedBlockedListListInvitesResponse) error {
	ownerChannel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}

	invites, err := model.BlockedListInvites(
		qm.Load(model.BlockedListInviteRels.BlockedList),
		qm.Load(model.BlockedListInviteRels.InviterChannel),
		model.BlockedListInviteWhere.InvitedChannelID.EQ(ownerChannel.ClaimID)).All(db.RO)

	var invitations []commentapi.SharedBlockedListInvitation
	for _, invite := range invites {
		if invite.R.BlockedList != nil && invite.R.InviterChannel != nil {
			list := commentapi.SharedBlockedList{}
			err = PopulateSharedBlockedList(&list, invite.R.BlockedList)
			if err != nil {
				return errors.Err(err)
			}

			if invite.R.BlockedList.InviteExpiration.Valid {
				expiresAt := invite.CreatedAt.Add(time.Duration(invite.R.BlockedList.InviteExpiration.Uint64) * time.Hour)
				if time.Now().After(expiresAt) {
					continue
				}
			}

			invitations = append(invitations, commentapi.SharedBlockedListInvitation{
				BlockedList: list,
				Invitation: commentapi.SharedBlockedListInvitedMember{
					InvitedByChannelName: invite.R.InviterChannel.Name,
					InvitedByChannelID:   invite.R.InviterChannel.ClaimID,
					InvitedChannelName:   ownerChannel.Name,
					InvitedChannelID:     ownerChannel.ClaimID,
					Status:               commentapi.InviteMemberStatusFrom(invite.Accepted, invite.CreatedAt, null.Uint64FromPtr(list.InviteExpiration)),
					InviteMessage:        invite.Message,
				},
			})
		}
	}

	reply.Invitations = invitations
	return nil
}

func invite(r *http.Request, args *commentapi.SharedBlockedListInviteArgs, reply *commentapi.SharedBlockedListInviteResponse) error {
	inviter, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}
	var blockedList *model.BlockedList
	if args.SharedBlockedListID != 0 {
		blockedList, err = model.BlockedLists(model.BlockedListWhere.ID.EQ(args.SharedBlockedListID)).One(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
	} else {
		blockedList, err = model.BlockedLists(model.BlockedListWhere.ChannelID.EQ(args.ChannelID)).One(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
	}

	if (blockedList.MemberInviteEnabled.Valid && !blockedList.MemberInviteEnabled.Bool) && inviter.ClaimID != blockedList.ChannelID {
		return api.StatusError{Err: errors.Err("shared blocked list %s does not have member inviting enabled", blockedList.Name)}
	}
	if !inviter.BlockedListInviteID.Valid {
		return api.StatusError{Err: errors.Err("channel %s is not authorized member of the shared blocked list %s", inviter.Name, blockedList.Name), Status: http.StatusUnauthorized}
	}
	if inviter.BlockedListInviteID.Uint64 != blockedList.ID {
		return api.StatusError{Err: errors.Err("channel %s is not a member of the shared blocked list %s", inviter.Name, blockedList.Name), Status: http.StatusBadRequest}
	}

	invitee, err := helper.FindOrCreateChannel(args.InviteeChannelID, args.InviteeChannelName)
	if err != nil {
		return errors.Err(err)
	}
	if invitee.BlockedListInviteID.Valid && invitee.BlockedListInviteID.Uint64 == blockedList.ID {
		return api.StatusError{Err: errors.Err("channel %s is already a member of the shared blocked list %s", invitee.Name, blockedList.Name), Status: http.StatusBadRequest}
	}
	where := model.BlockedListInviteWhere
	invite, err := model.BlockedListInvites(where.BlockedListID.EQ(args.SharedBlockedListID), where.InvitedChannelID.EQ(args.InviteeChannelID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if invite != nil {
		if !invite.Accepted.Valid {
			return api.StatusError{Err: errors.Err("channel %s already has an invite pending", invitee.Name), Status: http.StatusBadRequest}
		} else if !invite.Accepted.Bool {
			return api.StatusError{Err: errors.Err("channel %s already an invite and has rejected joining the shared blocked list %s", invitee.Name, blockedList.Name)}
		}
	}

	invite = &model.BlockedListInvite{
		BlockedListID:    blockedList.ID,
		InviterChannelID: inviter.ClaimID,
		InvitedChannelID: invitee.ClaimID,
		Message:          args.Message,
	}
	err = invite.Insert(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	return nil
}

func accept(r *http.Request, args *commentapi.SharedBlockedListInviteAcceptArgs, _ *commentapi.SharedBlockedListInviteResponse) error {
	channel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}

	where := model.BlockedListInviteWhere
	invite, err := model.BlockedListInvites(where.BlockedListID.EQ(args.SharedBlockedListID), where.InvitedChannelID.EQ(channel.ClaimID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	blockedList, err := model.BlockedLists(model.BlockedListWhere.ID.EQ(args.SharedBlockedListID)).One(db.RO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.StatusError{Err: errors.Err("there is no shared block list with id %d", args.SharedBlockedListID), Status: http.StatusBadRequest}
		}
		return errors.Err(err)
	}

	if blockedList == nil {
		return api.StatusError{Err: errors.Err("blocked list id %d does not exist", args.SharedBlockedListID), Status: http.StatusBadRequest}
	}

	if invite == nil {
		return api.StatusError{Err: errors.Err("channel %s does not have an invite for the shared block list %s to accept", args.ChannelName, blockedList.Name)}
	}

	if !args.Accepted {
		return rejectInvite(channel, invite)
	}

	return acceptInvite(channel, blockedList, invite)
}

func rejectInvite(channel *model.Channel, invite *model.BlockedListInvite) error {
	if channel.BlockedListID.IsZero() {
		return errors.Err("there is no blocked list currently contributing to reject an accepted invite")
	}
	if !invite.Accepted.Bool {
		return api.StatusError{Err: errors.Err("you have not accepted the invite yet")}
	}

	theirSBL, err := model.BlockedLists(model.BlockedListWhere.ChannelID.EQ(channel.ClaimID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	blockListID := null.Uint64{}
	if theirSBL != nil {
		blockListID = null.Uint64From(theirSBL.ID)
	}

	//Set blocks to their shared blocked list.
	blockedListCol := map[string]interface{}{model.BlockedEntryColumns.BlockedListID: blockListID}
	err = channel.CreatorChannelBlockedEntries().UpdateAll(db.RW, blockedListCol)
	if err != nil {
		return errors.Err(err)
	}

	channel.BlockedListID = blockListID
	channel.BlockedListInviteID = blockListID
	err = channel.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	invite.Accepted.SetValid(false)
	return errors.Err(invite.Update(db.RW, boil.Whitelist(model.BlockedListInviteColumns.Accepted)))
}

func acceptInvite(channel *model.Channel, blockedList *model.BlockedList, invite *model.BlockedListInvite) error {
	if blockedList.InviteExpiration.Valid {
		expiresAt := invite.CreatedAt.Add(time.Duration(blockedList.InviteExpiration.Uint64) * time.Hour)
		if time.Now().After(expiresAt) {
			return api.StatusError{Err: errors.Err("the invite expired at %s, and cannot be accepted", expiresAt.Format("2006-01-02 3:04:05 pm"))}
		}
	}
	blockedListID := null.Uint64From(blockedList.ID)
	acceptedInvites, err := channel.InvitedChannelBlockedListInvites(model.BlockedListInviteWhere.Accepted.EQ(null.BoolFrom(true))).All(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if len(acceptedInvites) > 0 {
		blockedListInviteCol := map[string]interface{}{model.BlockedListInviteColumns.Accepted: null.BoolFrom(false)}
		err := acceptedInvites.UpdateAll(db.RW, blockedListInviteCol)
		if err != nil {
			return errors.Err(err)
		}
	}
	invite.Accepted.SetValid(true)
	err = invite.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	blockedListCol := map[string]interface{}{model.BlockedEntryColumns.BlockedListID: blockedListID}
	err = channel.CreatorChannelBlockedEntries().UpdateAll(db.RW, blockedListCol)
	if err != nil {
		return errors.Err(err)
	}

	channel.BlockedListID = blockedListID
	channel.BlockedListInviteID = blockedListID
	err = channel.Update(db.RW, boil.Whitelist(model.ChannelColumns.BlockedListID, model.ChannelColumns.BlockedListInviteID))
	if err != nil {
		return errors.Err(err)
	}
	return nil
}

func rescind(r *http.Request, args *commentapi.SharedBlockedListRescindArgs, _ *commentapi.SharedBlockedListRescindResponse) error {
	ownerChannel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}

	var list *model.BlockedList

	list, err = model.BlockedLists(model.BlockedListWhere.ChannelID.EQ(ownerChannel.ClaimID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if list == nil {
		return api.StatusError{Err: errors.Err("blocked list not found"), Status: http.StatusNotFound}
	}

	invite, err := ownerChannel.InviterChannelBlockedListInvites(model.BlockedListInviteWhere.InvitedChannelID.EQ(args.InvitedChannelID), qm.Load(model.BlockedListInviteRels.InvitedChannel)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if invite == nil {
		return api.StatusError{Err: errors.Err("invite for %s not found", args.InvitedChannelName), Status: http.StatusBadRequest}
	}

	invitedChannel := invite.R.InvitedChannel
	invitedChannel.BlockedListInviteID = null.Uint64{}
	invitedChannel.BlockedListID = null.Uint64{}
	err = invitedChannel.Update(db.RW, boil.Whitelist(model.BlockedListInviteColumns.BlockedListID))
	if err != nil {
		return errors.Err(err)
	}

	err = invite.Delete(db.RW)
	if err != nil {
		return errors.Err(err)
	}

	return nil
}
