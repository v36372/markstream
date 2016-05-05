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
	"strconv"
	"sync"
	"time"
)

type Client struct {
	uuid string
	conn *websocket.Conn
	out  chan string
}

type Manager struct {
	clients map[string]*Client
	mutex   sync.Mutex
	// out     chan string
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
			c.out <- FloatToString(f64)
		}
		// m.out <- FloatToString(f64)
		log.Printf("sent output (%+v), sleeping for 1s...\n", f64)
		time.Sleep(time.Second)
	}
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
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
	cl.out = make(chan string)
	m.AddClient(cl)
	ws.Write([]byte("hehe1"))
	ws.Write([]byte("hehe2"))
	ws.Write([]byte("hehe3"))
	ws.Write([]byte("hehe4"))
	// go func() {
	// log.Print(FloatToString(<-cl.out))
	for val := range cl.out {
		// 	select {
		// 	// case <-c.Writer.CloseNotify():
		// 	// 	log.Printf("%s : disconnected\n", cl.uuid)
		// 	case out := <-cl.out:
		// 		// log.Print(FloatToString(<-cl.out))
		_, err := cl.conn.Write([]byte(val))
		if err != nil {
			m.DeleteClient(cl.uuid)
		}
		// 	case <-time.After(time.Second * 20):
		// 		log.Println("timed out")
		// 	default:
		// 		continue
		// 	}
	}
	// }()
}

func main() {
	m = NewManager()

	// m.out = make(chan string)
	go m.InitBackgroundTask()
	// go func() {
	// 	for {
	// 		select {
	// 		case out := <-m.out:
	// 			for _, c := range m.clients {
	// 				_, err := c.conn.Write([]byte(out))
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
