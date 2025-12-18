package handlers

import (
	"Means-to-an-End/message"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// message reader ..

func ReadMessage(conn net.Conn) (*message.Message, error) {
	var buf [9]byte
	_, err := io.ReadFull(conn, buf[:])
	if err != nil {
		return nil, err
	}

	msg := &message.Message{
		Type: buf[0],
	}
	copy(msg.A[:], buf[1:5])
	copy(msg.B[:], buf[5:9])

	return msg, nil
}

func WriteInt32(conn net.Conn, v int32) error {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v))
	_, err := conn.Write(buf[:])
	return err
}

// HandleConnection ...
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client Connected: %s\n", conn.RemoteAddr())
	store := make(map[int32]int32)

	for {
		msg, err := ReadMessage(conn)
		fmt.Println(msg)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client Disconnected: %s\n", conn.RemoteAddr())
			} else {
				log.Printf("Error: %v\n", err)
			}
			return
		}
		switch msg.Type {
		case 'I':
			ts := msg.IntA()
			price := msg.IntB()
			// fmt.Printf("INSERT: timestamp=%d price=%d\n", msg.A, msg.B)
			store[ts] = price

		case 'Q':
			min := msg.IntA()
			max := msg.IntB()

			if min > max {
				_ = WriteInt32(conn, 0)
				continue
			}

			var sum int64
			var count int64
			for ts, price := range store {
				if ts >= min && ts <= max {
					sum += int64(price)
					count++
				}
			}
			if count == 0 {
				_ = WriteInt32(conn, 0)
			} else {
				_ = WriteInt32(conn, int32(sum/count))
			}
		default:
			log.Println("Invalid Message type :", msg.Type)
			return
		}

	}
}
