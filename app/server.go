package main

import (
	"fmt"
	"io/ioutil"
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

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func do(conn net.Conn, dir string) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	check(err)
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
	} else if strings.Contains(path[1], "/files/") {
		filename := strings.Split(path[1], "/files/")[1]
		content, err := os.ReadFile(dir + "/" + filename)
		body := string(content)
		if os.IsNotExist(err) {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
			conn.Write([]byte(response))
		}

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}

func main() {

	dir := os.Args[2]

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

		go do(conn, dir)
	}
}
