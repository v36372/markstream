package main

import (
	"encoding/binary"
	// "fmt"
	// "github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"log"
	"math"
	"math/rand"
	// "net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	uuid string
	conn *websocket.Conn
	out  chan float64
}

type Manager struct {
	clients map[string]*Client
	mutex   sync.Mutex
}

var m *Manager

func NewManager() *Manager {
	m := new(Manager)
	m.clients = make(map[string]*Client)
	return m
}

func (m *Manager) AddClient(c *Client) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	log.Printf("add client: %s\n", c.uuid)
	m.clients[c.uuid] = c
}

func (m *Manager) DeleteClient(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	log.Println("delete client: %s", id)
	delete(m.clients, id)
}

func (m *Manager) InitBackgroundTask() {
	// log.Printf("aaaaaaaa")
	for {
		f64 := rand.Float64()
		log.Printf("active clients: %d\n", len(m.clients))
		for _, c := range m.clients {
			c.out <- f64
		}
		// m.out <- f64
		log.Printf("sent output (%+v), sleeping for 1s...\n", f64)
		time.Sleep(time.Second)
	}
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	// bits +=
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	// bytes = append(bytes, []byte(";")...)
	return bytes
}

// Echo the data received on the WebSocket.
func StreamServer(ws *websocket.Conn) {
	// io.Copy(ws, ws)
	cl := new(Client)
	cl.uuid = uuid.NewV4().String()
	cl.conn = ws
	cl.out = make(chan float64)
	m.AddClient(cl)
	go func() {
		for val := range cl.out {
			log.Printf("send")
			ws.Write(Float64bytes(val))
		}
	}()
}

func main() {
	m = NewManager()

	// m.out = make(chan float64)
	go m.InitBackgroundTask()
	// go func() {
	// 	for {
	// 		select {
	// 		case out := <-m.out:
	// 			for _, c := range m.clients {
	// 				_, err := c.conn.Write(Float64bytes(out))
	// 				if err != nil {
	// 					m.DeleteClient(c.uuid)
	// 				}
	// 			}
	// 		case <-time.After(time.Second * 20):
	// 			log.Println("timed out")
	// 		}
	// 	}
	// }()
	// ln, _ := net.Listen("tcp", ":8081") // accept connection on port
	http.Handle("/stream", websocket.Handler(StreamServer))

	// router := gin.Default()
	// router.LoadHTMLGlob("templates/*")
	// //router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl", gin.H{
	// 		"title": "Main website",
	// 	})
	// })

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	// go func() {
	// 	for {
	// 		conn, _ := ln.Accept() // run loop forever (or until ctrl-c)
	// 		cl := new(Client)
	// 		cl.uuid = uuid.NewV4().String()
	// 		cl.conn = conn
	// 		m.AddClient(cl)
	// 	}
	// }()

	// router.Run(":8080")
}
