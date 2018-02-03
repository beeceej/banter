package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/beeceej/banter"
	"github.com/beeceej/banter/pb"
)

var s1 = &banter.Server{Me: &pb.Peer{Name: "S1", Address: "localhost:8080", Port: ":8080"}, Quit: make(chan bool)}
var s2 = &banter.Server{Me: &pb.Peer{Name: "S2", Address: "localhost:8081", Port: ":8081"}, Quit: make(chan bool)}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() { s1.Register(123, wg) }()
	go func() { s2.Register(123, wg) }()

	a, b := s2.IssuePing(&pb.Peer{Address: "localhost:8080"}, time.Second*54)
	fmt.Println(a.String(), b)

	c, d := s1.IssuePing(s2.Me, time.Second*54)
	fmt.Println(c.String(), d)

	// s1.Quit <- true
	// s1.Quit <- true
	wg.Wait()
}
