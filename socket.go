package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var connections []*websocket.Conn
var mx sync.Mutex

type MuseWaveData struct {
	Wave  string
	Value float32
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Origin") != "http://"+r.Host {
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return
	// }

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	mx.Lock()
	connections = append(connections, conn)
	mx.Unlock()
}

func wsBroadcast(data MuseWaveData) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error broadcasting sensor data")
	}

	mx.Lock()
	for c := range connections {
		connections[c].WriteMessage(1, dataJson)
	}
	mx.Unlock()
}
