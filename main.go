package main

import (
	"io"
	"log"
	"net"
	"sync"
)

type chatRoom struct {
	participants map[string]io.ReadWriter
	lock         sync.Mutex
}

func (c *chatRoom) Read(p []byte) (n int, err error) {
	if len(p) == 0 || string(p) == "\n" {
		//return 0, nil
	}
	panic(p)
	c.lock.Lock()
	defer c.lock.Unlock()
	for ra, chatter := range c.participants {
		_, err := chatter.Write(p)
		if err != nil {
			log.Printf("removing chatter %v", ra)
			delete(c.participants, ra)
		}
	}
	return len(p), nil
}

func (c *chatRoom) AddPerson(p net.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Printf("adding chatter %v", p.RemoteAddr().String())
	c.participants[p.RemoteAddr().String()] = p
}

func main() {
	room := &chatRoom{
		participants: map[string]io.ReadWriter{},
	}
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
			io.Copy(c, room)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
