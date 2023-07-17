package main

import (
	"context"
	"io"
	"log"

	db "github.com/zezoamr/invoicer-service-grpc-Go/db"
	"github.com/zezoamr/invoicer-service-grpc-Go/proto/invoicer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type invoicerServer struct {
	invoicer.UnimplementedInvoicerServer
}

func (s *invoicerServer) SendVoiceMail(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.SendRequestStatus, error) {
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

func (s *invoicerServer) ReadUnSeenReceived(ctx context.Context, req *invoicer.ReadRequest) (*invoicer.MultipleinvoiceReply, error) {
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

func (s *invoicerServer) ReadAllReceived(req *invoicer.ReadRequest, stream invoicer.Invoicer_ReadAllReceivedServer) error {
	log.Printf("Got request for invoices for user : %v", req.Userid)

	items := db.DbReadAllReceived(dbconn, uint(req.Userid), int(req.Skip), int(req.Limit))
	for _, item := range items {
		res := &invoicer.InvoiceReply{
			ID:            uint32(item.ID),
			Audio:         item.Voice,
			Words:         item.Words,
			Wordsnotaudio: item.WordsNotAudio,
			FromID:        item.FromID,
			ToID:          item.ToID,
			Seen:          item.Seen,
		}

		// 2 second delay to simulate a long running process
		//time.Sleep(2 * time.Second)
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *invoicerServer) ReadAllSent(req *invoicer.ReadRequest, stream invoicer.Invoicer_ReadAllSentServer) error {
	log.Printf("Got request for invoices from user : %v", req.Userid)

	items := db.DbReadAllSent(dbconn, uint(req.Userid), int(req.Skip), int(req.Limit))
	for _, item := range items {
		res := &invoicer.InvoiceReply{
			ID:            uint32(item.ID),
			Audio:         item.Voice,
			Words:         item.Words,
			Wordsnotaudio: item.WordsNotAudio,
			FromID:        item.FromID,
			ToID:          item.ToID,
			Seen:          item.Seen,
		}

		// 2 second delay to simulate a long running process
		//time.Sleep(2 * time.Second)
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *invoicerServer) ReadMessageTime(stream invoicer.Invoicer_ReadMessageTimeServer) error {
	var messages []*invoicer.ReadMessageTime
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&invoicer.ReadMessageTimeResponse{SeenArray: messages})
		}
		if err != nil {
			return err
		}

		status, seen, time := db.DbReadMessageTime(dbconn, uint(req.Messageid))
		messages = append(messages, &invoicer.ReadMessageTime{Status: status, Read: seen, Readtime: timestamppb.New(time)})

	}
}

func (s *invoicerServer) MarkAsSeen(stream invoicer.Invoicer_MarkAsSeenServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		status, mid := db.DbMarkAsSeen(dbconn, uint(req.Messageid))

		if err := stream.Send(&invoicer.MarkRequestStatus{Status: status, Messageid: uint32(mid)}); err != nil {
			return err
		}
	}

}
