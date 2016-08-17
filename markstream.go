package markstream

import (
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"math"
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
	advertisement chan bool
}

func NewMarkStream() *MarkStream {
	ms := new(MarkStream)
	ms.userInputChan = make(chan string)
	ms.advertisement = make(chan bool)
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

func (ms *MarkStream) Process(fileName string, adsFileName string) {
	var audioData []float64
	var adsFile []float64
	audioData = ms.Read(fileName)
	adsFile = ms.Read(adsFileName)
	for {
		ms.Embedding(audioData)
		ms.advertisement <- true
		ms.Embedding(adsFile)
	}
}

func (ms *MarkStream) Input(adString string) {
	for {
		<-ms.advertisement
		time.Sleep(2 * time.Second)
		ms.userInputChan <- adString
	}
}
