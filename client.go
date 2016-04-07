package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// import "os"

func Float64frombytes(bytes []byte) float64 {
	fmt.Println("aaaaaaaaaaaaaa", bytes)
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func main() { // connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	for {
		// read in input from stdin
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Text to send: ")
		// text, _ := reader.ReadString('\n') // send to socket
		// fmt.Fprintf(conn, text+"\n")       // listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(Float64frombytes([]byte(message)))
	}
}
