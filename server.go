package banter

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/beeceej/banter/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Me   *pb.Peer
	Quit chan bool
}

func (s *Server) Send(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	switch msg.GetMsg() {
	case pb.Msg_PING:
		// fmt.Println(msg.GetOrigin().Address, "says", msg.GetMsg())
		r, e := s.IssuePong(msg.GetOrigin())
		// fmt.Printf("pong response %v %v", r.String(), e)
		return r, e
	case pb.Msg_PONG:
		// fmt.Println(msg.GetOrigin().Address, "says", msg.GetMsg())
		return &pb.Response{Status: pb.Status_OK}, nil
	}
	return nil, errors.New("unhandled msg")
}

func (s *Server) Dial(peer *pb.Peer) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		log.Printf("could not connect to | %v", err)
	}
	return conn, err
}

func (s *Server) Register(sessionID int, wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", s.Me.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	grpcsrv := grpc.NewServer()
	pb.RegisterBanterServer(grpcsrv, s)
	// Register reflection service on gRPC server.
	reflection.Register(grpcsrv)
	go func() {
		for {
			select {
			case b := <-s.Quit:
				if b {
					wg.Done()
					break
				}
			default:
			}
		}
	}()
	if err := grpcsrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
