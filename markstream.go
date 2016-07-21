package markstream

import (
	"bufio"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"log"
	"math"
	"os"
	"time"
)

const (
	MAG_THRES        = 0.0001
	SAMPLE_PER_FRAME = 22050
	BIN_PER_FRAME    = 800
	BIT_REPEAT       = 5
	PI               = math.Pi
)

type MarkStream struct {
	userInputChan chan string
	ConnManager   *Manager
}

func NewMarkStream() *MarkStream {
	ms := new(MarkStream)
	ms.userInputChan = make(chan string)
	ms.ConnManager = new(Manager)
	ms.ConnManager.clients = make(map[string]*Client)
	ms.ConnManager.audioDataChan = make(chan frame)

	return ms
}

// Echo the data received on the WebSocket.
func (ms *MarkStream) StreamServer(ws *websocket.Conn) {
	cl := new(Client)
	cl.uuid = uuid.NewV4().String()
	cl.conn = ws
	cl.exit = make(chan bool)
	ms.ConnManager.AddClient(cl)

	<-cl.exit
}

func (ms *MarkStream) Process(fileName string) {
	var l []float64
	l = ms.Read(fileName)

	ms.Embedding(l)
}

func (ms *MarkStream) Input() {
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Printf("Input your embedding string: ")
		text, _ := reader.ReadString('\n')
		ms.userInputChan <- text
		log.Println("Embedding...")
		time.Sleep(5 * time.Second)
	}
}
