package main

import (
	"github.com/v36372/markstream"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

type page struct {
	Title string
}

func main() {
	ms := markstream.NewMarkStream()
	fileName := string(os.Args[1])
	adsFileName := string(os.Args[2])
	adString := string(os.Args[3])

	go ms.Process(fileName, adsFileName)
	go ms.Input(adString)
	go ms.ConnManager.StreamToClients()

	http.Handle("/stream", websocket.Handler(ms.StreamServer))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
