package main

// import (
// 	"fmt"

// )

// func main() {

// 	// Output: bar
// }

import (
	"fmt"
	"net"
	"os"
	"time"

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
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
	}
	fmt.Println("Received message:", string(buffer))

	var item Item
	err = msgpack.Unmarshal(buffer, &item)
	if err != nil {
		panic(err)
	}
	fmt.Println(item)

	// write data to response
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("Your message is: %v from server. Received time: %v", string(buffer[:]), time)
	conn.Write([]byte(responseStr))

	// close conn
	conn.Close()
}
