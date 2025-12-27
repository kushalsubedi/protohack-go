package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

type Client struct {
	name string
	conn net.Conn
	send chan string
}


type Message struct {
	from string
	text string
}

var (
	joinCh      = make(chan *Client)
	leaveCh     = make(chan *Client)
	broadcastCh = make(chan Message)
)

func hub() {
	
	clients := make(map[*Client]bool)
	for {
		select {
		case c := <-joinCh:
			// Collect existing users
			var users []string
			for client := range clients {
				users = append(users, client.name)
			}

			if len(users) > 0 {
				c.send <- fmt.Sprintf("* The room contains: %s\n", strings.Join(users, ", "))
			}else{
				c.send <- fmt.Sprintf("* The room contains %s \n", strings.Join(users, ","))
			}

			// Add client to map after sending current users
			clients[c] = true

			// Notify other clients (excluding the new client)
			for client := range clients {
				if client != c {
					client.send <- fmt.Sprintf("* %s has entered the room\n", c.name)
				}
			}

		case c := <-leaveCh:
			delete(clients, c)
			close(c.send)
			for client := range clients {
				client.send <- fmt.Sprintf("* %s has left the room\n", c.name)
			}

		case msg := <-broadcastCh:
			for client := range clients {
				if client.name != msg.from {
					client.send <- fmt.Sprintf("[%s] %s", msg.from, msg.text)
				}
			}
		}
	}
}

func writeLoop(c *Client) {
	for msg := range c.send {
		// Write output to any Destination That implements io.write
		fmt.Fprint(c.conn, msg)
	}
}

func validName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	
	fmt.Println("Client was connected")
	_, err := conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	
	if err != nil {
		return
	}
	reader := bufio.NewReader(conn)

	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)

	if !validName(name) {
		fmt.Fprintln(conn, "Invalid name. Connection terminated.")

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.CloseWrite() // send FIN to client
		}
		return
	}
	client := &Client{
		name: name,
		conn: conn,
		send: make(chan string, 16),
	}
	go writeLoop(client)
	joinCh <- client

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			leaveCh <- client
			return
		}
		broadcastCh <- Message{
			from: name,
			text: msg,
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <PORT>:", os.Args[0])
		return
	}
	go hub()
	listenr, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		fmt.Printf("Couldnot Establish Connection :%s", err)
	}
	defer listenr.Close()
	for {
		conn, err := listenr.Accept()
		if err != nil {
			fmt.Printf("NO connection are active %s", err)
			continue
		}
		go HandleConn(conn)
	}
}
