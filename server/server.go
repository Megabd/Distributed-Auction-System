package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
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

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if server.auctionEnded {
			server.startAuction(scanner.Text())
		} else {
			fmt.Println("Auction already in progress")
		}

	}

}

func (s *Server) startAuction(text string) {

	if strings.Contains(text, "Start Auction") {
		s.auctionEnded = false
		s.maxBid = 0
		s.maxBidderId = 0
		fmt.Println("Auction Started")

		go func() {
			for i := 1; i < 240; i++ {
				if i == 239 {
					s.auctionEnded = true
				}
				time.Sleep(1 * time.Second)
			}
		}()

	} else {
		fmt.Println("Unknown Command, try: Start Auction")
	}

}

// given a bid, returns an outcome among {fail, success or exception}
func (s *Server) Bid(ctx context.Context, amount *proto.Amount) (*proto.Ack, error) {
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
