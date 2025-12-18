package main

import (
	"Means-to-an-End/handlers"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage %s:<port>", os.Args[0])
	}
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		log.Panic("Couldnot establish connection")
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic("Couldnt establish connection")
			return
		}

		go handlers.HandleConnection(conn)
	}
}
