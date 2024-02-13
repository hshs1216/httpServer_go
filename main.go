package main

import (
	"log"
	"net"
	"os"
)

const baseDir = "./" // サーバーの基準ディレクトリ

func main() {
	setupLogger()

	if len(os.Args) < 2 {
		log.Println("Usage: go run main.go <port>")
		os.Exit(1)
	}

	port := os.Args[1]

	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on :" + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting:", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
