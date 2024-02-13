package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Addr: %s, Connection\n", clientAddr)

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Addr: %s, Error reading request line: %v\n", clientAddr, err)
		return
	}

	method, path := parseRequestLine(requestLine)
	if method != "GET" {
		io.WriteString(conn, "HTTP/1.1 405 Method Not Allowed\r\n\r\n")
		log.Printf("Addr: %s, Method Not Allowed: %s\n", clientAddr, method)
		return
	}

	serveFile(conn, path, clientAddr)
}
