package banter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/beeceej/banter/pb"
	"google.golang.org/grpc"
)

func msgPing(origin *pb.Peer) *pb.Message {
	return &pb.Message{Msg: pb.Msg_PING, Origin: origin}
}

func msgPong(origin *pb.Peer) *pb.Message {
	return &pb.Message{Msg: pb.Msg_PONG, Origin: origin}
}

func (s *Server) IssuePing(peer *pb.Peer, deadline time.Duration) (*pb.Response, error) {
	var conn *grpc.ClientConn
	var err error
	if conn, err = s.Dial(peer); err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	c := pb.NewBanterClient(conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(deadline))
	defer cancel()
	// Contact the server and print out its response.
	if r, err := c.Send(ctx, msgPing(s.Me)); err != nil {

		return nil, err
	} else {
		fmt.Println(r.Status)
		return r, nil
	}
}

func (s *Server) IssuePong(peer *pb.Peer) (*pb.Response, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial(peer.GetAddress(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewBanterClient(conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	defer cancel()
	// Contact the server and print out its response.
	if r, err := c.Send(ctx, msgPong(s.Me)); err != nil {

		fmt.Println(err.Error())
		return nil, err
	} else {
		fmt.Println(r.Status)
		return r, nil
	}

}
