package banter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func init() {
	webportal = &portal{
		Shutdown: make(chan bool),
	}
}

type portal struct {
	Shutdown chan bool
}

var webportal *portal

func (s *Server) StartWebPortal(wg *sync.WaitGroup) chan bool {
	wg.Add(1)
	errCh := make(chan error)
	prtl := http.Server{Addr: ":3000"}

	go func() {
		for {
			select {
			case err := <-errCh:
				fmt.Println(err.Error())
				webportal.Shutdown <- true
				wg.Done()
			case <-webportal.Shutdown:
				fmt.Println("Shutdown")
				wg.Done()
			}
		}
	}()
	go func() {
		http.HandleFunc("/server", s.ServerInfo)
		http.HandleFunc("/shutdown", s.ShutdownPortal)

		if err := prtl.ListenAndServe(); err != nil {
			fmt.Println(err.Error())
			errCh <- err
		}
	}()

	return webportal.Shutdown
}

func (s *Server) ServerInfo(rw http.ResponseWriter, r *http.Request) {
	var (
		b   []byte
		err error
	)
	if b, err = json.MarshalIndent(s, "", " "); err != nil {
		fmt.Println(err.Error())
	}
	rw.Write(b)
}

func (s *Server) ShutdownPortal(rw http.ResponseWriter, r *http.Request) {
	webportal.Shutdown <- true
}
