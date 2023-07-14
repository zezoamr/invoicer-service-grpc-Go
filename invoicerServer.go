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
