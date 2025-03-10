package comments

import (
	"net/http"

	"github.com/lbryio/commentron/helper"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/sockety"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/sockety/socketyapi"

	"github.com/volatiletech/sqlboiler/boil"
)

func pin(_ *http.Request, args *commentapi.PinArgs) (commentapi.CommentItem, error) {
	var item commentapi.CommentItem
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).One(db.RO)
	if err != nil {
		return item, errors.Err(err)
	}

	claim, err := lbry.SDK.GetClaim(comment.LbryClaimID)
	if err != nil {
		return item, errors.Err(err)
	}
	if claim == nil {
		return item, errors.Err("could not resolve claim from comment")
	}

	channel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	claimChannel := claim.SigningChannel
	if claimChannel == nil {
		if claim.ValueType == "channel" {
			claimChannel = claim
		} else {
			return item, errors.Err("claim does not have a signing channel")
		}
	}

	err = lbry.ValidateSignatureFromClaim(claimChannel, args.Signature, args.SigningTS, args.CommentID)
	if err != nil {
		return item, err
	}
	comment.IsPinned = !args.Remove
	err = comment.Update(db.RW, boil.Infer())
	if err != nil {
		return item, errors.Err(err)
	}

	item = populateItem(comment, channel, 0)
	go sockety.SendNotification(socketyapi.SendNotificationArgs{
		Service: socketyapi.Commentron,
		Type:    "pinned",
		IDs:     []string{comment.LbryClaimID, "pins"},
		Data:    map[string]interface{}{"comment": item},
	})
	return item, nil
}
