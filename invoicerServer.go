package main

import (
	"context"

	"github.com/zezoamr/invoicer-service-grpc-Go/proto/invoicer"
)

type invoicerServer struct {
	invoicer.UnimplementedInvoicerServer
}

func (*invoicerServer) Create(context.Context, *invoicer.CreateRequest) (*invoicer.CreateResponseNoParmas, error) {
	return &invoicer.CreateResponseNoParmas{}, nil
}
