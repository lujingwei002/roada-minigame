package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/roada-go/gat/log"
)

type Type byte

const (
	Request  Type = 0x00
	Notify        = 0x01
	Response      = 0x02
	Push          = 0x03
)

const (
	msgRouteCompressMask = 0x01
	msgTypeMask          = 0x07
	msgRouteLengthMask   = 0xFF
	msgHeadLength        = 0x02
)

var types = map[Type]string{
	Request:  "Request",
	Notify:   "Notify",
	Response: "Response",
	Push:     "Push",
}

func (t Type) String() string {
	return types[t]
}

var (
	routes = make(map[string]uint16)
	codes  = make(map[uint16]string)
)

var (
	ErrWrongMessageType  = errors.New("wrong message type")
	ErrInvalidMessage    = errors.New("invalid message")
	ErrRouteInfoNotFound = errors.New("route info not found in dictionary")
	ErrWrongMessage      = errors.New("wrong message")
)

type Message struct {
	Type       Type
	ID         uint64
	Route      string
	Data       []byte
	compressed bool
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) String() string {
	return fmt.Sprintf("%s %s (%dbytes)", types[m.Type], m.Route, len(m.Data))
}

func (m *Message) Encode() ([]byte, error) {
	return Encode(m)
}

func routable(t Type) bool {
	return t == Request || t == Notify || t == Push
}

func invalidType(t Type) bool {
	return t < Request || t > Push

}

// ------------------------------------------
// |   type   |  flag  |       other        |
// |----------|--------|--------------------|
// | request  |----000-|<message id>|<route>|
// | notify   |----001-|<route>             |
// | response |----010-|<message id>        |
// | push     |----011-|<route>             |
// ------------------------------------------
func Encode(m *Message) ([]byte, error) {
	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}
	buf := make([]byte, 0)
	flag := byte(m.Type) << 1
	code, compressed := routes[m.Route]
	if compressed {
		flag |= msgRouteCompressMask
	}
	buf = append(buf, flag)
	if m.Type == Request || m.Type == Response {
		n := m.ID
		// variant length encode
		for {
			b := byte(n % 128)
			n >>= 7
			if n != 0 {
				buf = append(buf, b+128)
			} else {
				buf = append(buf, b)
				break
			}
		}
	}
	if routable(m.Type) {
		if compressed {
			buf = append(buf, byte((code>>8)&0xFF))
			buf = append(buf, byte(code&0xFF))
		} else {
			buf = append(buf, byte(len(m.Route)))
			buf = append(buf, []byte(m.Route)...)
		}
	}
	buf = append(buf, m.Data...)
	return buf, nil
}

func Decode(data []byte) (*Message, error) {
	if len(data) < msgHeadLength {
		return nil, ErrInvalidMessage
	}
	m := NewMessage()
	flag := data[0]
	offset := 1
	m.Type = Type((flag >> 1) & msgTypeMask)
	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}
	if m.Type == Request || m.Type == Response {
		id := uint64(0)
		// little end byte order
		// WARNING: must can be stored in 64 bits integer
		// variant length encode
		for i := offset; i < len(data); i++ {
			b := data[i]
			id += uint64(b&0x7F) << uint64(7*(i-offset))
			if b < 128 {
				offset = i + 1
				break
			}
		}
		m.ID = id
	}
	if offset >= len(data) {
		return nil, ErrWrongMessage
	}
	if routable(m.Type) {
		if flag&msgRouteCompressMask == 1 {
			m.compressed = true
			code := binary.BigEndian.Uint16(data[offset:(offset + 2)])
			route, ok := codes[code]
			if !ok {
				return nil, ErrRouteInfoNotFound
			}
			m.Route = route
			offset += 2
		} else {
			m.compressed = false
			rl := data[offset]
			offset++
			if offset+int(rl) > len(data) {
				return nil, ErrWrongMessage
			}
			m.Route = string(data[offset:(offset + int(rl))])
			offset += int(rl)
		}
	}
	if offset > len(data) {
		return nil, ErrWrongMessage
	}
	m.Data = data[offset:]
	return m, nil
}

func SetDictionary(dict map[string]uint16) {
	for route, code := range dict {
		r := strings.TrimSpace(route)
		// duplication check
		if _, ok := routes[r]; ok {
			log.Printf("duplicated route(route: %s, code: %d)\n", r, code)
		}
		if _, ok := codes[code]; ok {
			log.Printf("duplicated route(route: %s, code: %d)\n", r, code)
		}
		// update map, using last value when key duplicated
		routes[r] = code
		codes[code] = r
	}
}
