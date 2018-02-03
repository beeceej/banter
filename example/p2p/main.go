package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
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
	Address   string
	Port      string
}

func (s *server) Ping(ctx context.Context, in *pb.Client) (*pb.Response, error) {
	fmt.Println(ctx)
	conn, err := grpc.Dial(in.Address+in.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBanterClient(conn)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
	defer cancel()

	fmt.Printf("respondign to ping request from port %s", in.Port)
	r, err := c.Ping(ctx, &pb.Client{Address: in.Address, Port: in.Port})

	if err != nil {
		return nil, err
	}

	return &pb.Response{Status: r.GetStatus()}, nil
}

func registerServer(port string, sessionID int) {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println(port)
	pb.RegisterBanterServer(s, &server{Address: "127.0.0.1", Port: port, SessionID: sessionID})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registerClient(addresses []string, session int) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	// Set up a connection to the server.
	for _, address := range addresses {
		
			if err != nil {
				log.Printf("Could not ping %v | %s", address, err.Error())
				fmt.Println(&pb.Client{Address: strings.Split(address, ":")[0], Port: ":" + strings.Split(address, ":")[1]})

			} else {
				log.Printf(fmt.Sprintf("%s received message with status %s", address, r.GetStatus().String()))
				break
			}
		}
	}
}

func main() {
	port := os.Args[1:][0]
	peerAddresses := os.Args[2:]
	session := rand.Intn(1000000)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() { registerServer(port, session) }()
	registerClient(peerAddresses, session)
	wg.Wait()

}
