package main

import (
	"context"

	db "github.com/zezoamr/invoicer-service-grpc-Go/db"
	"github.com/zezoamr/invoicer-service-grpc-Go/proto/invoicer"
)

type invoicerServer struct {
	invoicer.UnimplementedInvoicerServer
}

func (*invoicerServer) SendVoiceMail(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.SendRequestStatus, error) {
	status, id := db.DbSendVoiceMail(dbconn, db.Invoice{
		Voice:         req.Message.Audio,
		Words:         req.Message.Words,
		WordsNotAudio: req.Message.Wordsnotaudio,
		FromID:        req.Userfrom,
		ToID:          req.Userto,
		Seen:          false,
	})
	return &invoicer.SendRequestStatus{Status: status, Messageid: uint32(id)}, nil
}

func (*invoicerServer) ReadUnSeenReceived(ctx context.Context, req *invoicer.ReadRequest) (*invoicer.MultipleinvoiceReply, error) {
	messages := db.DbReadUnSeenReceived(dbconn, uint(req.Userid), int(req.Skip), int(req.Limit))

	grpcItems := make([]*invoicer.InvoiceReply, len(messages))
	for i, item := range messages {
		grpcItems[i] = &invoicer.InvoiceReply{
			ID:            uint32(item.ID),
			Audio:         item.Voice,
			Words:         item.Words,
			Wordsnotaudio: item.WordsNotAudio,
			FromID:        item.FromID,
			ToID:          item.ToID,
			Seen:          item.Seen,
		}
	}

	return &invoicer.MultipleinvoiceReply{InvoiceReply: grpcItems}, nil
}
