package main

import (
	// "fmt"
	// b64 "encoding/base64"
	"encoding/binary"
	wavwriter "github.com/cryptix/wav"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	// "fmt"
	// "github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"log"
	// "math"
	"math/rand"
	// "net"
	"net/http"
	// "strconv"
	"bufio"
	"sync"
	"time"
)

var config struct {
	Header wav.Header
}

var input chan string

type frame []int16

const (
	MAG_THRES        = 0.0001
	SAMPLE_PER_FRAME = 22050
	BIN_PER_FRAME    = 800
	BIT_REPEAT       = 5
	PI               = math.Pi
)

type Client struct {
	uuid string
	conn *websocket.Conn
	out  chan frame
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
		// for _, c := range m.clients {
		// c.out <- FloatToString(f64)
		// }
		// m.out <- FloatToString(f64)
		log.Printf("sent output (%+v), sleeping for 1s...\n", f64)
		time.Sleep(time.Second)
	}
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 18, 64)
}

func Float64bytes(float float64) []byte {
	bits := math.Float32bits(float32(float))
	// bits +=
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	// bytes = append(bytes, []byte(";")...)
	return bytes
}

func Int16bytes(integer int16) []byte {
	bytes := make([]byte, 2)
	// big edian
	bytes[1] = uint8(integer >> 8)
	bytes[0] = uint8(integer & 0xff)
	return bytes
}

func Int16ArrayByte(f []int16) []byte {
	bytes := make([]byte, 0)
	for _, val := range f {
		bytes = append(bytes, Int16bytes(val)...)
	}
	// log.Println(len(bytes))
	return bytes
}

func FloatArrayByte(f []float64) []byte {
	bytes := make([]byte, 0)
	for _, val := range f {
		bytes = append(bytes, Float64bytes(val)...)
	}
	// log.Println(len(bytes))
	return bytes
}

// Echo the data received on the WebSocket.
func StreamServer(ws *websocket.Conn) {
	// io.Copy(ws, ws)
	cl := new(Client)
	cl.uuid = uuid.NewV4().String()
	cl.conn = ws
	cl.out = make(chan frame)
	m.AddClient(cl)
	// ws.Write([]byte("hehe1"))
	// ws.Write([]byte("hehe2"))
	// ws.Write([]byte("hehe3"))
	// ws.Write([]byte("hehe4"))
	// go func() {
	// log.Print(FloatToString(<-cl.out))
	for f := range cl.out {
		// 	select {
		// 	// case <-c.Writer.CloseNotify():
		// 	// 	log.Printf("%s : disconnected\n", cl.uuid)
		// 	case out := <-cl.out:
		// a := f
		// a[0] = 1
		// _, err := cl.conn.Write([]byte(b64.StdEncoding.EncodeToString(Int16ArrayByte(f))))
		err := websocket.Message.Send(cl.conn, Int16ArrayByte(f))
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

func Process() {
	// var filename = string(os.Args[1]) + ".wav"
	// var outputfile = string(os.Args[1]) + "_wm.wav"
	// var watermark = string(os.Args[2])
	var filename = "RWC_60s/RWC_002.wav"
	// var outputfile = "RWC_60s/RWC_001_wm.wav"
	// var watermark = "Nguyen Trong Tin"

	var l []float64
	l = Read(filename)

	// var mag []float64
	// var phs []float64
	// var currentpos int
	Embedding(l)

	// var newWav []float64
	// newWav = Reconstruct(mag, phs, l[currentpos:len(l)])

	// Write(newWav, outputfile)
}

func Input() {
	for {
		reader := bufio.NewReader(os.Stdin)
		log.Printf("Nhap zo di: ")
		text, _ := reader.ReadString('\n')
		input <- text
		// log.Printf("ahihi")
	}
}

func main() {
	m = NewManager()
	input = make(chan string)
	go Process()
	go Input()

	// m.out = make(chan string)
	// go m.InitBackgroundTask()

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
}

func Read(filename string) []float64 {
	file, _ := os.Open(filename)
	reader, _ := wav.New(file)

	l, _ := reader.ReadFloatsScale(reader.Samples)

	config.Header = reader.Header
	return l
}

func PrepareString(info string) string {
	var stringbit = ""
	byteArray := []byte(info)
	for _, char := range byteArray {
		n := int64(char)
		substr := strconv.FormatInt(n, 2)
		if len(substr) < 8 {
			length := len(substr)
			for j := 1; j <= 8-length; j++ {
				substr = "0" + substr
			}
		}
		stringbit += substr
	}

	return stringbit
}

func Embedding(l []float64) {
	// mag := make([]float64, 0)
	// phs := make([]float64, 0)

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	// var stringbit = PrepareString(watermark)
	// var stringbit = PrepareString("Nguyen Trong Tin")

	var pos = 0
	for i < len(l) {
		var subl = l[j : i+1]
		select {
		case watermark := <-input:
			log.Println("zo roi ne")
			submag := make([]float64, i+1-j)
			subphs := make([]float64, i+1-j)
			log.Println(watermark)
			var stringbit = PrepareString(watermark)
			for pos < len(stringbit) {
				subfourier := fft.FFTReal64(subl)
				var count = 0
				var bitrepeat = 0

				for k, x := range subfourier {
					submag[k], subphs[k] = cmplx.Polar(x)
					if submag[k] < MAG_THRES || k == 0 {
						continue
					}
					if pos < len(stringbit) && count < BIN_PER_FRAME {
						subphs[k] = QIMEncode(submag[k], subphs[k], int(stringbit[pos]))
						count++
						bitrepeat++
					}
					if bitrepeat == BIT_REPEAT {
						bitrepeat = 0
						pos++
					}
				}

				cmplxArray := make([]complex128, len(subl))
				for i, _ := range cmplxArray {
					cmplxArray[i] = cmplx.Rect(submag[i], subphs[i])
				}
				newWav := fft.IFFTRealOutput(cmplxArray)
				Wav16bit := Scale(newWav)
				for _, c := range m.clients {
					c.out <- Wav16bit
				}
			}
		default:
			// log.Println("gi z ta ?")
			Wav16bit := Scale(subl)
			log.Println(Wav16bit[0])
			for _, c := range m.clients {
				c.out <- Wav16bit
			}
		}
		j = i + 1
		i += SAMPLE_PER_FRAME
		if len(l)-i > 0 && len(l)-i < SAMPLE_PER_FRAME {
			i = len(l) - 1
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func ReconstructWithoutAppend(mag []float64, phs []float64) []float64 {
	cmplxArray := make([]complex128, len(mag))
	for i, _ := range mag {
		cmplxArray[i] = cmplx.Rect(mag[i], phs[i])
	}

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	var newWav = make([]float64, 0)
	for i < len(mag) {
		var subcmplx = cmplxArray[j : i+1]
		subIFFT := fft.IFFTRealOutput(subcmplx)
		newWav = append(newWav, subIFFT...)
		j = i + 1
		i += SAMPLE_PER_FRAME
		if len(mag)-i >= 0 && len(mag)-i < SAMPLE_PER_FRAME {
			i = len(mag) - 1
		}
	}

	return newWav
}

func Scale(wav []float64) []int16 {
	var newWav = make([]int16, len(wav))
	for i, x := range wav {
		integer := int16(x * math.MaxInt16)
		newWav[i] = integer
	}

	return newWav
}

func Reconstruct(mag []float64, phs []float64, original []float64) []float64 {
	cmplxArray := make([]complex128, len(mag))
	for i, _ := range mag {
		cmplxArray[i] = cmplx.Rect(mag[i], phs[i])
	}

	var i = SAMPLE_PER_FRAME - 1
	var j = 0
	var newWav = make([]float64, 0)
	for i < len(mag) {
		var subcmplx = cmplxArray[j : i+1]
		subIFFT := fft.IFFTRealOutput(subcmplx)
		newWav = append(newWav, subIFFT...)
		j = i + 1
		i += SAMPLE_PER_FRAME
		if len(mag)-i >= 0 && len(mag)-i < SAMPLE_PER_FRAME {
			i = len(mag) - 1
		}
	}

	newWav = append(newWav, original...)
	return newWav
}

func Write(newWav []float64, outputfile string) {
	wavOut, _ := os.Create(outputfile)
	defer wavOut.Close()

	meta := wavwriter.File{
		Channels:        1,
		SampleRate:      config.Header.SampleRate,
		SignificantBits: config.Header.BitsPerSample,
	}

	writer, _ := meta.NewWriter(wavOut)
	defer writer.Close()

	for n := 0; n < len(newWav); n += 1 {
		integer := int16(newWav[n] * math.MaxInt16)
		writer.WriteInt16(integer)
	}
}

func QIMEncode(mag float64, phs float64, bit int) float64 {
	step := [5]float64{PI / 18, PI / 14, PI / 10, PI / 6, PI / 2}
	var stepsize = findStep(mag)
	if bit == 48 {
		return math.Floor(phs/step[stepsize]+0.5) * step[stepsize]
	} else {
		return math.Floor(phs/step[stepsize])*step[stepsize] + step[stepsize]/2
	}
}

func findStep(mag float64) int32 {
	var sMag = mag / (0.005)
	var group = math.Ceil(sMag / 0.2)
	if group == 0 {
		group = 0
	}
	if group > 4 {
		group = 4
	}
	return int32(group)
}
