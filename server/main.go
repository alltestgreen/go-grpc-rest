package main

import (
	"context"
	"hello/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

// Runs a gRPC server using Protobuf serialisation. SubmitOrder functionality computes the total amount of the order.
func main() {

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterOrderServiceServer(srv, &server{})

	log.Println("gRPC server started")

	if e := srv.Serve(listener); e != nil {
		panic(err)
	}
}

func (s *server) SubmitOrder(ctx context.Context, order *proto.Order) (*proto.Confirmation, error) {
	log.Println("Processing order:", order.GetId())

	total := order.GetProduct().GetPrice() * float32(order.GetQuantity())

	return &proto.Confirmation{Amount: total}, nil
}
