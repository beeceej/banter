package banter

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/beeceej/banter/pb"
)

var peers = []*pb.Peer{
	&pb.Peer{Name: "S1", Address: "localhost", Port: "8080"},
	&pb.Peer{Name: "S2", Address: "localhost", Port: "8081"},
	&pb.Peer{Name: "S3", Address: "localhost", Port: "8082"},
	&pb.Peer{Name: "S4", Address: "localhost", Port: "8083"},
	&pb.Peer{Name: "S5", Address: "localhost", Port: "8084"},
	&pb.Peer{Name: "S6", Address: "localhost", Port: "8085"},
	&pb.Peer{Name: "S7", Address: "localhost", Port: "8086"},
	&pb.Peer{Name: "S8", Address: "localhost", Port: "8087"},
	&pb.Peer{Name: "S9", Address: "localhost", Port: "8088"},
	&pb.Peer{Name: "S10", Address: "localhost", Port: "8089"},
	&pb.Peer{Name: "S11", Address: "localhost", Port: "8090"},
	&pb.Peer{Name: "S12", Address: "localhost", Port: "8091"},
}

var servers = []*Server{
	&Server{
		Me: &pb.Peer{
			Name:    "S1",
			Address: "localhost",
			Port:    "8080",
		},
		Peers: peers[1:],
		Quit:  make(chan bool),
	},
	&Server{Me: peers[1], Quit: make(chan bool)},
	&Server{Me: peers[2], Quit: make(chan bool)},
	&Server{Me: peers[3], Quit: make(chan bool)},
	&Server{Me: peers[4], Quit: make(chan bool)},
	&Server{Me: peers[5], Quit: make(chan bool)},
	&Server{Me: peers[6], Quit: make(chan bool)},
	&Server{Me: peers[7], Quit: make(chan bool)},
	&Server{Me: peers[8], Quit: make(chan bool)},
	&Server{Me: peers[9], Quit: make(chan bool)},
	&Server{Me: peers[10], Quit: make(chan bool)},
	&Server{Me: peers[11], Quit: make(chan bool)},
}

func TestBroadcast(t *testing.T) {
	wg := new(sync.WaitGroup)
	s1 := servers[0]
	go func() { s1.Register(rand.Intn(1000), wg) }()
	for _, v := range servers[1:] {
		go func(s *Server) { s.Register(rand.Intn(1000), wg) }(v)
	}
	time.Sleep(time.Millisecond * 50)

	results := s1.Broadcast(MsgPing(s1.Me), time.Second*2)

	for k, v := range results {
		fmt.Println(k, v)
	}

	for _, v := range servers[1:] {
		v.Quit <- true
	}
	s1.Quit <- true

	wg.Wait()
}
