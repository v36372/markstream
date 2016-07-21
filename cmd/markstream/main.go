package main

import (
	"golang.org/x/net/websocket"
	"$GOPATH/src/markstream"
	"net/http"
)

type page struct {
	Title string
}

func main() {
	ms := markstream.NewMarkStream()
	go ms.Process()
	go ms.Input()
	go ms.ConnManager.StreamToClients()

	http.Handle("/stream", websocket.Handler(ms.StreamServer))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
