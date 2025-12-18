package message

import (
	"encoding/binary"
	"fmt"
)

// Message

type Message struct {
	Type byte
	A    [4]byte
	B    [4]byte
}

func (m *Message) IntA() int32 {
	return int32(binary.BigEndian.Uint32(m.A[:]))
}

func (m *Message) IntB() int32 {
	return int32(binary.BigEndian.Uint32(m.B[:]))
}

func (m *Message) String() string {
	return fmt.Sprintf(
		"Message{Type:%c, A:%d, B:%d}",
		m.Type,
		m.IntA(),
		m.IntB(),
	)
}
