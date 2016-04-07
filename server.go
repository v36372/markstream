package main

import (
	// "bufio"
	"fmt"
	// "github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net"
	// "net/http"
	// "strings" // only needed below for sample processing
	"encoding/binary"
	"math"
	"sync"
	"time"
)

type Client struct {
	uuid string
	conn net.Conn
}

type Manager struct {
	clients map[string]*Client
	out     chan float64
	mutex   sync.Mutex
}

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
	// log.Println("delete client: %s", c.uuid)
	delete(m.clients, id)
}

func (m *Manager) InitBackgroundTask() {
	for {
		f64 := rand.Float64()
		// log.Printf("active clients: %d\n", len(m.clients))
		// for _, c := range m.clients {
		// 	c.out <- f64
		// }
		m.out <- f64
		log.Printf("sent output (%+v), sleeping for 10s...\n", f64)
		time.Sleep(time.Second)
	}
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	// bits +=
	bytes := make([]byte, 16)
	binary.LittleEndian.PutUint64(bytes, bits)
	bytes = append(bytes, []byte("\n")...)
	return bytes
}

func main() {
	fmt.Println("Launching server...")
	m := NewManager()
	m.out = make(chan float64)
	go m.InitBackgroundTask()
	ln, _ := net.Listen("tcp", ":8081") // accept connection on port
	// conn, _ := ln.Accept()              // run loop forever (or until ctrl-c)

	go func() {
		for {
			conn, _ := ln.Accept() // run loop forever (or until ctrl-c)
			cl := new(Client)
			cl.uuid = uuid.NewV4().String()
			// cl.out = make(chan float64)
			cl.conn = conn
			// defer m.DeleteClient(cl.uuid)
			m.AddClient(cl)
		}
	}()

	for { // will listen for message to process ending in newline (\n)
		// conn, _ := ln.Accept() // run loop forever (or until ctrl-c)
		// message, _ := bufio.NewReader(conn).ReadString('\n') // output message received
		// fmt.Print("Message Received:", string(message))      // sample process for string received
		// newmessage := "hello" // send new string back to client
		select {
		// case <-c.Writer.CloseNotify():
		// 	log.Printf("%s : disconnected\n", cl.uuid)
		case out := <-m.out:
			for _, c := range m.clients {
				c.conn.Write(Float64bytes(out))
			}
		case <-time.After(time.Second * 20):
			log.Println("timed out")
		}
	}
}
