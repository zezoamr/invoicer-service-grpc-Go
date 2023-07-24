package main

import (
	"log"

	inv "github.com/zezoamr/invoicer-service-grpc-Go/proto/invoicer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8089"
)

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := inv.NewInvoicerClient(conn)

	mySlice := []byte{'a', 'b'}
	words := []string{"hello", "world"}

	callSendVoiceMail(client, 1, 2, mySlice, words, true)
	callReadUnSeenReceived(client, 2, -1, -1)
	callReadAllReceived(client, 2, -1, -1)
	callReadAllSent(client, 1, -1, -1)
	callMarkAsSeen(client, []uint32{1})
	callReadMessageTime(client, []uint32{1})

}
