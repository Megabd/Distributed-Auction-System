package main

import (
	"Distributed-Auction-System/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
		return
	}

	defer conn.Close()

	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ctx := context.Background()
	client := proto.NewAuctionClient(conn)

	c := &clientStruct{
		id:  arg1,
		ctx: ctx,
	}

	fmt.Println("Started succesfully")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		c.bidOrInfo(client, scanner.Text())
	}
}

type clientStruct struct {
	id  int64
	ctx context.Context
}

func (c *clientStruct) bidOrInfo(client proto.AuctionClient, text string) {
	command := strings.Split(text, " ")

	if command[0] == "bid" {
		amount, err := strconv.ParseInt(command[1], 10, 32)
		if err != nil {
			log.Fatalf("invalid bid, try again.")
			return
		}

		c.bid(client, amount)

	} else if command[0] == "info" {
		c.info(client)

	} else {
		fmt.Println("Unknown command, try again")
	}
}

func (c *clientStruct) bid(client proto.AuctionClient, amount int64) {
	response, err := client.Bid(c.ctx, &proto.Amount{Id: c.id, Value: amount})
	if err != nil {
		log.Fatalf("Something went wrong")
		return
	}
	if response.Success {
		fmt.Println("Successful bid, you are now highest bidder")
	} else if !response.Success {
		fmt.Println("Your bid did not exceed the current highest bid/there is no active Auction")
	}
}

func (c *clientStruct) info(client proto.AuctionClient) {
	response, err := client.Result(c.ctx, &proto.Void{})
	if err != nil {
		log.Fatalf("Something went wrong")
		return
	}
	if response.Over {
		fmt.Println("Auction finished.")
	} else {
		fmt.Println("Auction ongoing.")
	}
	fmt.Printf("Highest bid: %v by bidder nr. %v.", response.Value, response.Id)
	fmt.Println()
}
