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

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("ClientLogs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ctx := context.Background()

	c := &clientStruct{
		id:      arg1,
		ctx:     ctx,
		servers: make(map[int32]proto.AuctionClient),
	}

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)
		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		server := proto.NewAuctionClient(conn)
		c.servers[int32(i)] = server

	}

	fmt.Println("Started succesfully")
	log.Printf("Client with ID: %v started succesfully", c.id)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var results [3]string
		for i, server := range c.servers {
			results[i] = c.bidOrInfo(server, scanner.Text())
		}
		for i := range results {
			if results[i] == results[0] && i != 0 || results[i] == results[1] && i != 1 || results[i] == results[2] && i != 2 {
				fmt.Println(results[i])
				break
			}
		}
	}
}

type clientStruct struct {
	id      int64
	ctx     context.Context
	servers map[int32]proto.AuctionClient
}

func (c *clientStruct) bidOrInfo(server proto.AuctionClient, text string) string {
	command := strings.Split(text, " ")

	if command[0] == "bid" {
		amount, err := strconv.ParseInt(command[1], 10, 32)
		if err != nil {
			log.Println("invalid bid, try again.")
			return "invalid bid, try again."
		}

		return c.bid(server, amount)

	} else if command[0] == "info" {
		return c.info(server)

	} else {
		return "Unknown command, try again"
	}
}

func (c *clientStruct) bid(server proto.AuctionClient, amount int64) string {
	response, err := server.Bid(c.ctx, &proto.Amount{Id: c.id, Value: amount})

	if err != nil {
		log.Printf("Something went wrong with bid from ID: %v \n", c.id)
		return "Something went wrong"
	}
	if response.Success {
		log.Printf("%v is now the highest bidder with bid: %v\n", c.id, amount)
		return "Successful bid, you are now highest bidder"
	} else {
		log.Printf("%v's bid did not exceed the current highest bid/there is no active Auction \n", c.id)
		return "Your bid did not exceed the current highest bid/there is no active Auction"
	}

}

func (c *clientStruct) info(server proto.AuctionClient) string {
	response, err := server.Result(c.ctx, &proto.Void{})

	var result string

	if err != nil {
		result = "Something went wrong\n"
		log.Printf("Something went wrong with auction info from/to ID: %v", c.id)
		return result
	}
	if response.Over {
		result = "Auction finished."

	} else {
		result = "Auction ongoing.\n"
	}
	result2 := "Highest bid: " + strconv.Itoa(int(response.Value)) + " by bidder nr." + strconv.Itoa(int(response.Id))
	result = result + "\n" + result2
	return result

}
