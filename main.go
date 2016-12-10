package main

import (
	"bufio"
	"log"
	"net"
	"sync"
	"time"
)

type chatRoom struct {
	participants map[string]net.Conn
	lock         sync.Mutex
	chatter      [][]byte
}

func (c *chatRoom) ReadAll() {
	for _, conn := range c.participants {
		go func() {
			reader := bufio.NewReader(conn)
			b, err := reader.ReadBytes('\n')
			if err != nil {
				panic(err)
			}
			c.writeAll(conn.RemoteAddr().String(), b)
		}()
	}
}

func (c *chatRoom) writeAll(sender string, p []byte) {
	for _, conn := range c.participants {
		if conn.RemoteAddr().String() == sender {
			continue
		}
		go conn.Write(p)
	}
}

func (c *chatRoom) RunRoom() {
	for {
		c.ReadAll()
		time.Sleep(1e9)
	}
}

func (c *chatRoom) AddPerson(p net.Conn) {
	log.Printf("adding chatter %v", p.RemoteAddr().String())
	c.participants[p.RemoteAddr().String()] = p
}

func main() {
	room := &chatRoom{
		participants: map[string]net.Conn{},
	}
	go room.RunRoom()
	// Listen on TCP port 2000 on all interfaces.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Add to chat room
			room.AddPerson(c)
			for {
				time.Sleep(10e9)
			}
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
