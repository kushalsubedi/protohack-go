package main

import (
	"fmt"
	"json-tcp/handeler"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <port>", os.Args[0])
	}
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	fmt.Println("server running on port :", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection")
			continue
		}
		go handeler.HandleConnection(conn)
	}
}
