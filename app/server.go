package main

import (
	"fmt"
	"log"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

func get_header(s string) (string, string) {
	line := strings.Split(s, ":")
	return line[0], line[1][1:]
}

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
	} else if strings.Contains(path[1], "/echo/") {
		ans := strings.Split(path[1], "/echo/")[1]
		reponse := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(ans), ans)
		conn.Write([]byte(reponse))
	} else if strings.Contains(path[1], "/user-agent") {
		_, v := get_header(lines[2])
		reponse := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(v), v)
		conn.Write([]byte(reponse))
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
