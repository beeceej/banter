package main

import (
	"math/rand"
	"sync"

	"github.com/beeceej/banter"
	"github.com/beeceej/banter/pb"
)

var s = &banter.Server{
	Me: &pb.Peer{
		Name:    "S1",
		Address: "localhost",
		Port:    "8080",
	},
	Peers: []*pb.Peer{},
	Quit:  make(chan bool),
}

func main() {
	wg := new(sync.WaitGroup)
	go func() { s.Register(rand.Intn(1000), wg) }()
	s.StartWebPortal(wg)
	wg.Wait()
}
