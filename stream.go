package markstream

import (
	// "bufio"
	// "github.com/kataras/iris"
	// "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"log"
	// "net/http"
	// "os"
	"sync"
	// "time"
)

type Client struct {
	uuid string
	conn *websocket.Conn
	exit chan bool
}

type frame []int16

type Manager struct {
	clients       map[string]*Client
	mutex         sync.Mutex
	audioDataChan chan frame
}

func NewManager() *Manager {
	m := new(Manager)
	m.clients = make(map[string]*Client)
	return m
}

func (m *Manager) StreamToClients() {
	for audioData := range m.audioDataChan {
		for _, cl := range m.clients {
			err := websocket.Message.Send(cl.conn, Int16ArrayByte(audioData))
			if err != nil {
				cl.exit <- true
				m.DeleteClient(cl.uuid)
			}
		}
	}
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
