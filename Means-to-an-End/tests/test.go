package main

import (
	"encoding/binary"
	"log"
	"net"
)

func sendInsert(conn net.Conn, ts int32, price int32) {
	buf := make([]byte, 9)
	buf[0] = 'I'
	binary.BigEndian.PutUint32(buf[1:5], uint32(ts))
	binary.BigEndian.PutUint32(buf[5:9], uint32(price))
	conn.Write(buf)
}

func sendQuery(conn net.Conn, min int32, max int32) int32 {
	buf := make([]byte, 9)
	buf[0] = 'Q'
	binary.BigEndian.PutUint32(buf[1:5], uint32(min))
	binary.BigEndian.PutUint32(buf[5:9], uint32(max))
	conn.Write(buf)

	resp := make([]byte, 4)
	conn.Read(resp)
	return int32(binary.BigEndian.Uint32(resp))
}

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	sendInsert(conn, 1000, 10)
	sendInsert(conn, 2000, 20)
	sendInsert(conn, 3000, 30)

	mean := sendQuery(conn, 1000, 3000)
	log.Println("Mean:", mean) // should print 20
}
