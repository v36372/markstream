package main

import (
	"github.com/v36372/markstream"
	"golang.org/x/net/websocket"
	"net/http"
)

type page struct {
	Title string
}

func main() {
	ms := markstream.NewMarkStream()
	fileName := string(os.Args[1])

	go ms.Process(fileName)
	go ms.Input()
	go ms.ConnManager.StreamToClients()

	http.Handle("/stream", websocket.Handler(ms.StreamServer))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
