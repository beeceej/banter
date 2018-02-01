package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/beeceej/banter/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type server struct {
	SessionID int
}

func (s *server) Ping(ctx context.Context, in *pb.Address) (*pb.Response, error) {
	status := func() pb.Status {
		success := rand.Intn(5) > 2
		if success {
			return pb.Status_OK
		}
		return pb.Status_ERROR
	}()

	return &pb.Response{Status: status}, nil
}

func registerServer(port string, sessionID int) {
	// should be env variable

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBanterServer(s, &server{SessionID: sessionID})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registerClient(addresses []string, session int) {
	// Set up a connection to the server.
	for _, address := range addresses {
		conn, err := grpc.Dial(fmt.Sprintf(address), grpc.WithInsecure())

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewBanterClient(conn)

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
		defer cancel()
		// Contact the server and print out its response.

		for times := 0; times < 3; times++ {
			r, err := c.Ping(ctx, &pb.Address{Address: address})
			if err != nil {
				log.Printf("Could not ping %v", address)
			} else {
				log.Printf(fmt.Sprintf("%s received message with status %s", address, r.GetStatus().String()))
				break
			}
		}
	}
}

func main() {
	port := os.Args[1:][0]
	peerAddresses := os.Args[1:]
	session := rand.Intn(1000000)
	go func() { registerServer(port, session) }()
	for {
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		registerClient(peerAddresses, session)
	}
}
