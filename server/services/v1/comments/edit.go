package comments

import (
	"net/http"
	"time"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/boil"
)

func edit(args *commentapi.EditArgs) (*commentapi.CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).One(db.RO)
	if err != nil {
		return nil, errors.Err(err)
	}
	if comment == nil {
		return nil, api.StatusError{Err: errors.Err("could not find comment with id %s", args.CommentID), Status: http.StatusBadRequest}
	}
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(comment.ChannelID.String)).One(db.RO)
	if err != nil {
		return nil, errors.Err(err)
	}
	if channel == nil {
		return nil, api.StatusError{Err: errors.Err("channel id %s could not be found"), Status: http.StatusBadRequest}
	}
	err = lbry.ValidateSignatureAndTS(comment.ChannelID.String, args.Signature, args.SigningTS, args.Comment)
	if err != nil {
		return nil, err
	}

	comment.Body = args.Comment
	comment.IsPinned = false
	comment.Signature.SetValid(args.Signature)
	comment.Signingts.SetValid(args.SigningTS)
	comment.Timestamp = int(time.Now().Unix())
	err = comment.Update(db.RW, boil.Infer())
	if err != nil {
		return nil, errors.Err(err)
	}
	item := populateItem(comment, channel, 0)
	return &item, nil
}
