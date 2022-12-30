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

	request "github.com/lucasfth/disys-exam-active-replication/grpc"
	"google.golang.org/grpc"
)

func main(){
	var id int32
	log.Printf("Enter id below:")
	fmt.Scanln(&id)

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
		table: make(map[int32]int32),
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

	resp := &request.PutResponse{}
	resp.Response = true
	srv.Send(resp);	
	return nil;
}

func (s *server) SendBid(in *request.Put, srv request.BiddingService_SendBidServer) error{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	resp := &request.PutResponse{}

	var output bytes.Buffer
	output.WriteString(fmt.Sprintf("Update hash %v with %v :", in.Hash, in.Val))

	if in.Hash < 0 {
		resp.Response = false
		output.WriteString(" fail")
	} else {
		s.table[in.Hash] = in.Val
		resp.Response = true
		output.WriteString(" success")
	}

	log.Print(output.String())
	
	srv.Send(resp)
	return nil
}

func (s *server) RequestCurrentResult(in *request.Get, srv request.BiddingService_RequestCurrentResultServer) error {
	val := s.table[in.Hash]
	
	log.Printf("Hash %v	with %v", in.Hash, val)
	
	resp := &request.GetResponse{}
	resp.Val = val

	srv.Send(resp);
	return nil;
}

type server struct{
	mutex 			sync.Mutex
	ownPort 		int32
	table 			map[int32]int32
	ctx 			context.Context
	request.UnimplementedBiddingServiceServer
}