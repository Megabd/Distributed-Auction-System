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
		auctionEnded: false,
		maxBid:       0,
		maxBidderId:  0,
	}

	grpcServer := grpc.NewServer()

	auction.RegisterAuctionServer(grpcServer, server) //Registers the server to the gRPC server.

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
	fmt.Printf("Server started successfully")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		server.startAuction(scanner.Text())

	}

}

func (s *Server) startAuction(text string) {

	if strings.Contains(text, "Start Auction") {
		s.auctionEnded = false
		s.maxBid = 0
		s.maxBidderId = 0

		go func() {
			for i := 1; i < 240; i++ {
				if i == 239 {
					s.auctionEnded = true
				}
				time.Sleep(1 * time.Second)
			}
		}()

	}

}

// given a bid, returns an outcome among {fail, success or exception}
func (s *Server) Bid(ctx context.Context, amount *proto.Amount) (*proto.Ack, error) {
	//register bidder if new
	ack := &proto.Ack{}

	if amount.Value <= 0 { //or amount lower than previous bid from bidder
		ack.Success = false
	} else if s.maxBid < amount.Value {
		s.maxBid = amount.Value
		s.maxBidderId = amount.Id
		ack.Success = true
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
