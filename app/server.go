package main

import (
	"fmt"
	"log"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func do(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err.Error())
	}
	res := string(buf)
	lines := strings.Split(res, "\r\n")
	path := strings.Split(lines[0], " ")
	if path[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go do(conn)
	}
}
