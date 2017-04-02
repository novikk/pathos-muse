package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/hypebeast/go-osc/osc"
)

var sensorStatus = [4]int{4, 4, 4, 4}

func getWaveHandler(wave string) func(msg *osc.Message) {
	return func(msg *osc.Message) {
		//fmt.Println(msg)

		total := float32(0)
		qty := 0
		for i := range msg.Arguments {
			if !math.IsNaN(float64(msg.Arguments[i].(float32))) {
				total += msg.Arguments[i].(float32)
				qty++
			}
		}

		val := total / float32(qty)

		wsBroadcast(MuseWaveData{
			Wave:  wave,
			Value: val,
		})

		//fmt.Println(wave, "Value:", val)
	}
}

func getSensorStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, sensorStatus)
}

func main() {
	addr := "127.0.0.1:5000"
	server := &osc.Server{Addr: addr}

	server.Handle("/muse/elements/experimental/mellow", func(msg *osc.Message) {
		wsBroadcast(MuseWaveData{
			Wave:  "Mellow",
			Value: msg.Arguments[0].(float32) / float32(len(msg.Arguments)),
		})

		//fmt.Println("Stress:", msg.Arguments[0].(float32)/float32(len(msg.Arguments)))
	})

	server.Handle("/muse/elements/alpha_absolute", getWaveHandler("Alpha"))
	server.Handle("/muse/elements/beta_absolute", getWaveHandler("Beta"))
	server.Handle("/muse/elements/gamma_absolute", getWaveHandler("Gamma"))
	server.Handle("/muse/elements/theta_absolute", getWaveHandler("Theta"))

	server.Handle("/muse/elements/horseshoe", func(msg *osc.Message) {
		//fmt.Println("Status:", msg.Arguments)
		for i := range msg.Arguments {
			sensorStatus[i] = int(msg.Arguments[i].(float32))
		}
	})

	go func() {
		fmt.Println(server.ListenAndServe())
	}()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/getSensorStatus", getSensorStatus)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}
