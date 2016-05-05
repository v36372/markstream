package main

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	uuid string
	out  chan float64
}

type Manager struct {
	clients map[string]*Client
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
		for _, c := range m.clients {
			c.out <- f64
		}
		log.Printf("sent output (%+v), sleeping for 10s...\n", f64)
		time.Sleep(time.Second)
	}
}

func main() {
	r := gin.Default()
	m := NewManager()

	go m.InitBackgroundTask()

	r.GET("/", func(c *gin.Context) {
		cl := new(Client)
		cl.uuid = uuid.NewV4().String()
		cl.out = make(chan float64)

		defer m.DeleteClient(cl.uuid)
		m.AddClient(cl)

		select {
		case <-c.Writer.CloseNotify():
			log.Printf("%s : disconnected\n", cl.uuid)
		case out := <-cl.out:
			log.Printf("%s : received %+v\n", out)
			c.JSON(http.StatusOK, gin.H{
				"output": out,
			})
		case <-time.After(time.Second * 20):
			log.Println("timed out")
		}
	})

	r.Run()
}
