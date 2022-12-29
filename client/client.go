package main

// This template draws inspiration from the following source:
// github.com/lucasfth/go-ass5
// Which was created by chbl, fefa and luha
// The template itself was created by chbl and luha

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	request "github.com/lucasfth/disys-exam-active-replication/grpc"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(0)
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	c := &client{downedServerInt: 0}

	log.Printf("Enter username below:")
	fmt.Scanln(&c.name)

	// To change log location, outcomment below

	// path := fmt.Sprintf("clientlog_%s", c.name)
	// f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	log.Printf("Welcome %s", c.name)

	c.downedServers = make (map[int32]bool) // init map of downed servers

	// Connect to all servers
	for i := 0; i < 3; i++ { // Will iterate through ports 5001, 5002, 5003
		dialNum  := int32(5001 + i)
		dialNumString := fmt.Sprintf(":%v", dialNum) 

		conn, err := grpc.Dial(dialNumString, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		
		// create stream
		client := request.NewBiddingServiceClient(conn)
		in := &request.ClientHandshake{ClientPort: dialNum, Name: c.name} 
		//bidStream, err := client.SendBid(context.Background(), )
		stream, err := client.Handshake(ctx, in)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Connected to server %v and responded %s", dialNum, resp)
		}
		if err != nil {
			log.Fatalf("Cannot receive %v", err)
		}
		c.servers = append(c.servers, client)
		c.downedServers[dialNum] = false
		time.Sleep(4 * time.Second)
	}

	c.communication()
}

func (c *client) communication() {
	for { // Communication loop
		rand.Seed(time.Now().UnixNano()) // ensure "random" number is different each time
		delay := int32(rand.Intn(4))
		time.Sleep(time.Duration(delay) * time.Second)
		c.sendBids()
	}
}

func (c *client) sendBids(){
	responses := make([]int32, len(c.servers))
	for i := 0; i < len(c.servers); i++ { // Send bid to all servers
		response, _ := c.sendBid(int32(i))
		responses[i] = response
	}
	logicResponse := c.logic(responses)
	
	if logicResponse == -1 {
		log.Print("Tried to bid but found EXCEPTION")
		os.Exit(1)
	}

	log.Printf("---------Inc was %v", logicResponse)
}

func (c *client) sendBid(iteration int32) (int32, error) {
	in := &request.Bid{Name: c.name}
	stream, err := c.servers[iteration].SendBid(context.Background(), in)
	if err != nil {
		serverDown(iteration, c)
		return -1, err
	}
	resp, err := stream.Recv()
	return resp.GetResponse(), err
}

func serverDown (iteration int32, c *client) (bool){
	if (!c.downedServers[5001 + iteration]) {
		log.Printf("Server %v is down", (5001 + iteration))
		c.downedServers[5001 + iteration] = true
		return true // Server has just crashed
	}
	return false // Server was already down
}

func auctionFinished(winnerName string, winnerAmount int32, name string) {
	if (winnerName == name) {
		log.Printf("Won the auction with bid %v", winnerAmount)
	} else {
		log.Printf("Lost the auction")
	}
	os.Exit(1)
}

func (c *client) logic(responses []int32) (int32) {
	for i := 0; i < len(responses); i++ {
		// log.Printf("Response was: %s ,on i: %v", responses[i], i)
		if responses[i] > 0 {
			// log.Printf("Went into success")
			return responses[i]
		} else if (i == len(responses) - 1) {
			return -1
		}
	}
	return -1
}

type client struct {
	name 			string
	downedServers 	map[int32]bool
	downedServerInt int32
	servers 		[]request.BiddingServiceClient
}
