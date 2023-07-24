package main

import (
	"context"
	"io"
	"log"
	"time"

	inv "github.com/zezoamr/invoicer-service-grpc-Go/proto/invoicer"
)

func callSendVoiceMail(client inv.InvoicerClient, userfrom uint32, userto uint32, aud []byte, wordsarr []string, wordsbool bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SendVoiceMail(ctx, &inv.CreateRequest{
		Userfrom: userfrom,
		Userto:   userto,
		Message: &inv.Invoice{
			Audio:         aud,
			Words:         wordsarr,
			Wordsnotaudio: wordsbool,
		},
	})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("success status: %t , if true created voice mail with id: %b", res.Status, res.Messageid)
}

func callReadUnSeenReceived(client inv.InvoicerClient, userid int, skip int32, limit int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.ReadUnSeenReceived(ctx, &inv.ReadRequest{Userid: uint32(userid), Skip: skip, Limit: limit})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("recivied that amount: %d", len(res.GetInvoiceReply()))
}

func callReadAllReceived(client inv.InvoicerClient, userid int, skip int32, limit int32) {
	log.Printf("Streaming started")
	stream, err := client.ReadAllReceived(context.Background(), &inv.ReadRequest{Userid: uint32(userid), Skip: skip, Limit: limit})
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}
		log.Println(message.FromID, message.ToID, message.Audio[0:9], message.Words[0:9], message.Wordsnotaudio)
	}

	log.Printf("Streaming finished")
}

func callReadAllSent(client inv.InvoicerClient, userid int, skip int32, limit int32) {
	log.Printf("Streaming started")
	stream, err := client.ReadAllSent(context.Background(), &inv.ReadRequest{Userid: uint32(userid), Skip: skip, Limit: limit})
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}
		log.Println(message.FromID, message.ToID, message.Audio[0:9], message.Words[0:9], message.Wordsnotaudio)
	}

	log.Printf("Streaming finished")
}

func callMarkAsSeen(client inv.InvoicerClient, markids []uint32) {
	log.Printf("Bidirectional Streaming started")
	stream, err := client.MarkAsSeen(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	waitc := make(chan struct{}) //synchronize the execution of the goroutine with the main function

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message.Messageid, message.Status)
		}
		close(waitc)
	}()

	for _, id := range markids {
		req := &inv.MarkRequest{
			Messageid: id,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		//time.Sleep(2 * time.Second)
	}

	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional Streaming finished")
}

func callReadMessageTime(client inv.InvoicerClient, mIds []uint32) {

	stream, err := client.ReadMessageTime(context.Background())
	if err != nil {
		log.Fatalf("Could not send ids : %v", err)
	}

	for _, id := range mIds {
		req := &inv.ReadMessageTimeRequest{
			Messageid: id,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		log.Printf("Sent request with id: %d", id)
		//time.Sleep(2 * time.Second) for delay simu
	}

	res, err := stream.CloseAndRecv()
	log.Printf("Client Streaming finished")
	if err != nil {
		log.Fatalf("Error while receiving %v", err)
	}
	log.Printf("recivied that amount: %d", len(res.GetSeenArray()))
}
