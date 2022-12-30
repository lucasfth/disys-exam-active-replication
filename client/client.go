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
		actionType := int32(rand.Intn(2)) // 0 = bid, 1 = request
		randomHash := int32(rand.Intn(1000))

		if actionType == 0 {
			// log.Printf("Action type: Bid")
			randomVal := int32(rand.Intn(1000))
			c.sendBids(randomHash, randomVal)
			time.Sleep(4 * time.Second)
		} else {
			// log.Printf("Action type: Request")
			c.requestCurrentResults(randomHash)
			time.Sleep(4 * time.Second)
		}
	}
}

func (c *client) sendBids(hash int32, val int32) {
	responses := make([]bool, len(c.servers))
	for i := 0; i < len(c.servers); i++ { // Send bid to all servers
		response, _ := c.sendBid(int32(i), hash, val)
		responses[i] = response
	}
	resp, logicResponse := c.logic(responses)
	
	if logicResponse == "Exception" {
		log.Print("Tried to bid but found EXCEPTION")
		os.Exit(1)
	}

	log.Printf("---------Hash %v with val %v : %t", hash, val, resp)
}

func (c *client) sendBid(iteration int32, hash int32, val int32) (bool, error) {
	in := &request.Put{Hash: hash, Val: val}
	stream, err := c.servers[iteration].SendBid(context.Background(), in)
	if err != nil {
		serverDown(iteration, c)
		return false, err
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

func (c *client) requestCurrentResults(hash int32) (currentRelaventBid int32){
	var value int32 

	for i := 0; i < len(c.servers); i++ { // Request current result from all servers
		resp, err := c.requestCurrentResult(int32(i), hash)
		if err != nil {
			if serverDown(int32(i), c) {
				c.downedServerInt++
			}
			continue
		}
		if (c.downedServerInt == int32(len(c.servers))){ // All servers are probably down
			log.Printf("Tried to request but found EXCEPTION")
			os.Exit(1)
		}
		
		value = resp.Val
	}
	log.Printf("---------Current val on hash %v is %v", hash, value)
	return value
}

func (c *client) requestCurrentResult(iteration int32, hash int32)(*request.GetResponse, error){
	in := &request.Get{Hash: hash}
	stream, err := c.servers[iteration].RequestCurrentResult(context.Background(), in)
	if err != nil {
		return nil, err
	}
	resp, err := stream.Recv()
	return resp, err
}

func (c *client) logic(responses []bool) (bool, string) {
	for i := 0; i < len(responses); i++ {
		// log.Printf("Response was: %s ,on i: %v", responses[i], i)
		if responses[i] {
			// log.Printf("Went into success")
			return true, "Success"
		} else if responses[i] {
			// log.Printf("Went into fail")
			return false, "Fail"
		} else if (i == len(responses) - 1) {
			return false, "Exception"
		}
	}
	return false, "Exception"
}

type client struct {
	name 			string
	downedServers 	map[int32]bool
	downedServerInt int32
	servers 		[]request.BiddingServiceClient
}
