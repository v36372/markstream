package markstream

import (
	"bufio"
	// "github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"time"
	"math"
)

type MarkStream struct {
  userInputChan chan string
  log *log.Logger
  connManager *Manager
}

const (
  MAG_THRES        = 0.0001
  SAMPLE_PER_FRAME = 22050
  BIN_PER_FRAME    = 800
  BIT_REPEAT       = 5
  PI               = math.Pi
)

// Echo the data received on the WebSocket.
func (ms *MarkStream) StreamServer(ws *websocket.Conn) {
	cl := new(Client)
	cl.uuid = uuid.NewV4().String()
	cl.conn = ws
	// cl.out = make(chan frame)
	ms.connManager.AddClient(cl)
	// for f := range cl.out {
	// 	err := websocket.Message.Send(cl.conn, Int16ArrayByte(f))
	// 	if err != nil {
	// 		m.DeleteClient(cl.uuid)
	// 	}
	// }
}

func (ms *MarkStream) Process() {
	var filename = "RWC_60s/RWC_002.wav"

	var l []float64
	l = ms.Read(filename)

	ms.Embedding(l)
}

func (ms *MarkStream) Input() {
	for {
		reader := bufio.NewReader(os.Stdin)
		log.Printf("Input your embedding string: ")
		text, _ := reader.ReadString('\n')
		ms.userInputChan <- text
		// m.embedd <- "start"
		log.Println("Embedding...")
		time.Sleep(5 * time.Second)
	}
}
