package main

import (
	"sync"

	"github.com/hypebeast/go-osc/osc"
)

func main() {
	addr := "127.0.0.1:5000"
	server := &osc.Server{Addr: addr}

	server.Handle("/muse/eeg", func(msg *osc.Message) {
		osc.PrintMessage(msg)
	})

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		server.ListenAndServe()
		wg.Done()
	}()

	wg.Wait()
}
