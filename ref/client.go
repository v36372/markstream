package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

func Float64frombytes(bytes []byte) float64 {
	fmt.Println("aaaaaaaaaaaaaa", bytes[:8])
	bits := binary.LittleEndian.Uint64(bytes[:8])
	float := math.Float64frombits(bits)
	return float
}

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	defer conn.Close()
	for {
		message, _ := bufio.NewReader(conn).ReadString(';')
		fmt.Println(Float64frombytes([]byte(message)))
	}
}
