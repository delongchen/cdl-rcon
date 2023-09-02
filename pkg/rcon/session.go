package rcon

import (
	"fmt"
	"log"
	"net"
	"time"
)

const PING_ID int32 = 0

type RCONSession struct {
	conn net.Conn
	addr string
	In   chan *BasicPacket
	Out  chan *BasicPacket
	ID   int32
}

type AuthError struct{}

func (*AuthError) Error() string {
	return "auth not ok"
}

func NewSession(addr string) *RCONSession {
	return &RCONSession{
		ID:   1,
		conn: nil,
		addr: addr,
		In:   make(chan *BasicPacket),
		Out:  make(chan *BasicPacket),
	}
}

func (s *RCONSession) dial() error {
	conn, err := net.Dial("tcp", s.addr)
	if err != nil {
		return err
	}

	s.conn = conn

	return nil
}

func (s *RCONSession) readLoop() {
	buf := make([]byte, 4096)

	for {
		_, err := s.conn.Read(buf)

		if err != nil {
			close(s.Out)
			return
		}

		s.Out <- FromBytes(buf)
	}
}

func (s *RCONSession) sendLoop() {
	for {
		p := <-s.In
		_, err := s.conn.Write(p.ToBytes())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *RCONSession) startPingLoop() {
	time.Sleep(time.Second)
	for {
		ping := &BasicPacket{
			ID:   PING_ID,
			Type: SERVERDATA_RESPONSE_VALUE,
			Body: "",
		}
		s.In <- ping
		time.Sleep(10 * time.Second)
	}
}

func (s *RCONSession) StartLoop() error {
	if s.conn == nil {
		err := s.dial()
		if err != nil {
			return err
		}
		fmt.Printf("dial %s ok\n", s.addr)
	}

	go s.sendLoop()
	go s.readLoop()

	s.In <- &BasicPacket{
		ID:   PING_ID,
		Type: SERVERDATA_AUTH,
		Body: "4789516729Chen",
	}

	authResult := <-s.Out
	if authResult.ID != PING_ID {
		close(s.In)
		close(s.Out)
		return &AuthError{}
	}

	go s.startPingLoop()

	return nil
}
