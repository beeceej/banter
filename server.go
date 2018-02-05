package banter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/beeceej/banter/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Server struct {
	Me    *pb.Peer
	Peers []*pb.Peer
	Quit  chan bool `json:"-"`
}

func (s *Server) Send(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	switch msg.GetMsg() {
	case pb.Msg_PING:
		return s.IssuePong(msg.GetOrigin(), time.Second*1)
	case pb.Msg_PONG:
		return &pb.Response{Status: pb.Status_OK}, nil
	}
	return nil, errors.New("unhandled msg")
}

func (s *Server) Dial(peer *pb.Peer) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%s", peer.GetAddress(), peer.GetPort()), grpc.WithInsecure())
	if err != nil {
		log.Printf("could not connect to | %v", err)
	}
	return conn, err
}

func (s *Server) Register(sessionID int, wg *sync.WaitGroup) {
	wg.Add(1)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Me.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	grpcsrv := grpc.NewServer()
	pb.RegisterBanterServer(grpcsrv, s)
	reflection.Register(grpcsrv)
	go func() {
		for {
			select {
			case b := <-s.Quit:
				if b {
					fmt.Println("done")
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

func (s *Server) IssuePing(peer *pb.Peer, deadline time.Duration) (r *pb.Response, err error) {
	var conn *grpc.ClientConn
	if conn, err = s.Dial(peer); err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBanterClient(conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(deadline))
	defer cancel()
	if r, err = c.Send(ctx, MsgPing(s.Me)); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Server) IssuePong(peer *pb.Peer, deadline time.Duration) (r *pb.Response, err error) {
	var conn *grpc.ClientConn
	conn, err = s.Dial(peer)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBanterClient(conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(deadline))
	defer cancel()
	if r, err = c.Send(ctx, MsgPing(s.Me)); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Server) Broadcast(msg *pb.Message, deadline time.Duration) map[string]pb.Status {
	ch := make(chan *pb.Response)
	results := make(map[string]pb.Status)
	for _, v := range s.Peers {
		go func(p *pb.Peer) {
			r, err := s.Issue(p, msg, deadline)
			if err != nil {
				ch <- &pb.Response{Status: pb.Status_ERROR}
			}
			ch <- r
		}(v)
	}
	for _, peer := range s.Peers {
		select {
		case r := <-ch:
			results[fmt.Sprintf("%s:%s", peer.GetAddress(), peer.GetPort())] = r.GetStatus()
		}
	}
	return results
}

func (s *Server) Issue(peer *pb.Peer, msg *pb.Message, deadline time.Duration) (r *pb.Response, err error) {
	switch msg.GetMsg() {
	case pb.Msg_PING:
		r, err = s.IssuePing(peer, deadline)
	case pb.Msg_PONG:
		r, err = s.IssuePong(peer, deadline)
	}
	return r, err
}
