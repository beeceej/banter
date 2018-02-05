package banter

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/beeceej/banter/pb"
)

func TestBroadcast(t *testing.T) {
	wg := new(sync.WaitGroup)
	s1 := &Server{
		Me: &pb.Peer{
			Name:    "S1",
			Address: "localhost",
			Port:    "8080",
		},
		Peers: []*pb.Peer{
			&pb.Peer{Name: "S2", Address: "localhost", Port: "8081"},
			&pb.Peer{Name: "S3", Address: "localhost", Port: "8082"},
			&pb.Peer{Name: "S4", Address: "localhost", Port: "8083"},
		},
		Quit: make(chan bool),
	}
	var servers []*Server
	var srv *Server
	for i, p := range s1.Peers {
		if i != 1 {
			srv = &Server{
				Me: &pb.Peer{
					Name:    p.GetName(),
					Address: p.GetAddress(),
					Port:    p.GetPort(),
				},
				Quit: make(chan bool),
			}
			servers = append(servers, srv)
			go func(s *Server) {
				s.Register(rand.Intn(1000), wg)
			}(srv)
		}
	}
	go func() { s1.Register(rand.Intn(1000), wg) }()
	time.Sleep(time.Millisecond * 1110)
	results := s1.Broadcast(MsgPing(s1.Me), time.Second*3)
	fmt.Println(results)
	shutdown(servers)
	s1.Quit <- true
	wg.Wait()
	if results["localhost:8081"] != pb.Status_ERROR {
		t.Fail()
	}
}

func shutdown(servers []*Server) {
	for _, v := range servers {
		v.Quit <- true
	}
}
