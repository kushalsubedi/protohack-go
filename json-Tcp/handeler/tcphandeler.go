package handeler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"json-tcp/messages"
	"json-tcp/util"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	encoder := json.NewEncoder(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}
		fmt.Printf("RaW: %s", line)

		var msg messages.Request
		if err := json.Unmarshal(line, &msg); err != nil {
			encoder.Encode(messages.ErrorResponse{Error: "malformed request"})
			return
		}

		if msg.Method != "isPrime" || msg.Number == nil {
			encoder.Encode(messages.ErrorResponse{Error: "malformed request"})
			return
		}
		n := *msg.Number

		resp := messages.PrimeResponse{
			Method: "isPrime",
			Prime:  util.CheckPrime(n),
		}
		if err := encoder.Encode(resp); err != nil {
			return
		}
	}
}
