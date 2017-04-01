package main

import (
	"fmt"
	"net/http"

	"github.com/hypebeast/go-osc/osc"
)

func main() {
	addr := "127.0.0.1:5000"
	server := &osc.Server{Addr: addr}

	server.Handle("/muse/elements/experimental/mellow", func(msg *osc.Message) {
		//fmt.Println(msg)

		total := float32(0)
		for i := range msg.Arguments {
			total += msg.Arguments[i].(float32)
		}

		wsBroadcast(MuseWaveData{
			Wave:  "EEG",
			Value: total / float32(len(msg.Arguments)),
		})

		fmt.Println("Value:", total/float32(len(msg.Arguments)))
	})

	server.Handle("/muse/elements/horseshoe", func(msg *osc.Message) {
		fmt.Println("Status:", msg.Arguments)
	})

	go func() {
		fmt.Println(server.ListenAndServe())
	}()

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}
