package main

import (
	"net"
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

type Item struct {
	Foo    string
	Params map[string]interface{}
}

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	b, err := msgpack.Marshal(&Item{Foo: "kaka", Params: map[string]interface{}{"a": 1, "b": 2}})
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(b)
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 2048)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println("Received message:", string(received))

	conn.Close()
}
