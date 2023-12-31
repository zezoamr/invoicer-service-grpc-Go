package main

import (
	"log"
	"net"

	db "github.com/zezoamr/invoicer-service-grpc-Go/db"
	"github.com/zezoamr/invoicer-service-grpc-Go/proto/invoicer"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB
)

func handleError(err error, str string) {
	if err != nil {
		log.Fatalf("error has happened %s %v", str, err)
	}
}

func main() {

	port := ":8089"
	lis, err := net.Listen("tcp", port)
	handleError(err, "cannot create listner on port"+port)

	serverRegister := grpc.NewServer()
	service := &invoicerServer{}

	dbconn = db.InitDatabase()

	invoicer.RegisterInvoicerServer(serverRegister, service)
	err = serverRegister.Serve(lis)
	handleError(err, "not able to serve")
}
