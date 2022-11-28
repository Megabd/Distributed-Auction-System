package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"Distributed-Auction-System/proto"
	auction "Distributed-Auction-System/proto"

	"google.golang.org/grpc"
)

type Server struct {
	auction.UnimplementedAuctionServer
	auctionEnded bool
	maxBid       int64
	maxBidderId  int64
}

func main() {

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	ownPort := 5000

	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))

	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := &Server{
		auctionEnded: true,
		maxBid:       0,
		maxBidderId:  0,
	}

	grpcServer := grpc.NewServer()
	auction.RegisterAuctionServer(grpcServer, server) //Registers the server to the gRPC server.

	fmt.Println("Server started successfully")
	log.Printf("Server with port: %v started succesfully", ownPort)
	main2(ownPort + 1)
	main2(ownPort + 2)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func main2(portId int) {

	ownPort := portId

	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))

	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := &Server{
		auctionEnded: true,
		maxBid:       0,
		maxBidderId:  0,
	}

	grpcServer := grpc.NewServer()
	auction.RegisterAuctionServer(grpcServer, server) //Registers the server to the gRPC server.

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	fmt.Println("Server started successfully")
	log.Printf("Server with port: %v started succesfully", ownPort)
}

func (s *Server) startAuction() {

	s.auctionEnded = false
	s.maxBid = 0
	s.maxBidderId = 0

	go func() {
		for i := 1; i < 60; i++ {
			if i == 59 {
				s.auctionEnded = true
				log.Printf("The winner is: %v with a max bid of: %v\n", s.maxBidderId, s.maxBid)
			}
			time.Sleep(1 * time.Second)
		}
	}()

}

// given a bid, returns an outcome among {fail, success or exception}
func (s *Server) Bid(ctx context.Context, amount *proto.Amount) (*proto.Ack, error) {
	if s.auctionEnded {
		s.startAuction()
	}

	ack := &proto.Ack{}
	if !s.auctionEnded {
		if amount.Value <= 0 {
			ack.Success = false
		} else if s.maxBid < amount.Value {
			s.maxBid = amount.Value
			s.maxBidderId = amount.Id
			ack.Success = true
		} else {
			ack.Success = false
		}
	} else {
		ack.Success = false
	}
	return ack, nil
}

// if the auction is over, it returns the result, else highest bid.
func (s *Server) Result(ctx context.Context, void *proto.Void) (*proto.Outcome, error) {
	reply := &proto.Outcome{Id: s.maxBidderId, Value: s.maxBid, Over: s.auctionEnded}
	return reply, nil
}
