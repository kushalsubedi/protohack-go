package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}

	addr := ":" + os.Args[1]

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer ln.Close()

	fmt.Println("Echo server running on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handle(conn) // handle clients concurrently
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	// io.Copy automatically:
	// - Reads until EOF
	// - Writes each chunk exactly as received
	// - Works with raw binary data
	// - Efficient (uses large buffers)
	_, err := io.Copy(conn, conn)
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println("Connection error:", err)
	}

	// When io.Copy returns EOF:
	// - client has closed its write side
	// - server has echoed everything
	// Now closing conn tells the client we are done.
}
