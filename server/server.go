package main
// This template draws inspiration from the following source:
// github.com/lucasfth/go-ass5
// Which was created by chbl, fefa and luha
// The template itself was created by chbl and luha

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	request "github.com/lucasfth/disys-exam-active-replication/grpc"
	"google.golang.org/grpc"
)

func main(){
	var id int32
	log.Printf("Enter id below:")
	fmt.Scanln(&id)

	var endHour, endMin int
	log.Printf("Enter auction time below in hour and min:")
	fmt.Scanln(&endHour, &endMin)
	end := time.Date (time.Now().Year(), time.Now().Month(), time.Now().Day(), endHour, endMin, 0, 0, time.Local)

	// To change log location, outcomment below

	// path := fmt.Sprintf("serverlog_%v", id)
	// f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	log.Printf("Hello %v", id)
	
	// Create port
	port := int32(5000 + id)
	portString := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &server{
		mutex: sync.Mutex{},
		ownPort: port,
		currentBid: 0,
		currentBidOwner: "", 
		isOver: false, 
		auctionEnd: end,
		ctx: context.Background(),
	}

	// create grpc server
	s := grpc.NewServer()
	request.RegisterBiddingServiceServer(s, server)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func(s *server) Handshake(in *request.ClientHandshake, srv request.BiddingService_HandshakeServer) error {
	log.Printf("Handshake 	%s", in.Name)

	resp := &request.BidResponse{}
	resp.Response = 0
	srv.Send(resp);	
	return nil;
}

func (s *server) SendBid(in *request.Bid, srv request.BiddingService_SendBidServer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	resp := &request.BidResponse{}

	if (time.Until(s.auctionEnd) <= 0) {
		if !s.isOver { log.Printf("--- Auction is over, %s won with bid %v ---", s.currentBidOwner, s.currentBid) }
		s.isOver = true
	}

	var output bytes.Buffer
	output.WriteString(fmt.Sprintf("Bid\t\t%s", in.Name))

	if s.isOver {
		resp.Response = -1
		output.WriteString(fmt.Sprintf("\t%v with TODO but auction over, winner: %s , with: %v", resp.Response, s.currentBidOwner, s.currentBid))
	} else {
		s.currentBid += 1
		s.currentBidOwner = in.Name
		resp.Response = s.currentBid
		output.WriteString(fmt.Sprintf("\t%v with increment", resp.Response))
	}

	log.Print(output.String())
	
	srv.Send(resp)
	return nil
}

type server struct{
	mutex 			sync.Mutex
	ownPort 		int32
	currentBid 		int32
	currentBidOwner	string
	isOver 			bool
	auctionEnd 		time.Time
	ctx 			context.Context
	request.UnimplementedBiddingServiceServer
}