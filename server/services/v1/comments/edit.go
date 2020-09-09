package comments

import (
	"net/http"
	"time"

	"github.com/lbryio/commentron/commentapi"

	"github.com/lbryio/lbry.go/extras/api"

	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/volatiletech/sqlboiler/boil"
)

func edit(args *commentapi.EditArgs) (*commentapi.CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	if comment == nil {
		return nil, api.StatusError{Err: errors.Err("could not find comment with id %s", args.CommentID), Status: http.StatusBadRequest}
	}
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(comment.ChannelID.String)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	if channel == nil {
		return nil, api.StatusError{Err: errors.Err("channel id %s could not be found"), Status: http.StatusBadRequest}
	}
	err = lbry.ValidateSignature(comment.ChannelID.String, args.Signature, args.SigningTS, args.Comment)
	if err != nil {
		return nil, err
	}

	comment.Body = args.Comment
	comment.Signature.SetValid(args.Signature)
	comment.Signingts.SetValid(args.SigningTS)
	comment.Timestamp = int(time.Now().Unix())
	err = comment.UpdateG(boil.Infer())
	if err != nil {
		return nil, errors.Err(err)
	}
	item := populateItem(comment, channel)
	return &item, nil
}
