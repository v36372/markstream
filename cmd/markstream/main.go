package main

import (
	"markstream"
  // "net/http"
  // "golang.org/x/net/websocket"
  "github.com/kataras/iris"
)

func main() {
	ms := markstream.NewMarkStream()
	// go m.InitBackgroundTask()
	go ms.Process()
	go ms.Input()

	http.Handle("/stream", websocket.Handler(ms.StreamServer))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
