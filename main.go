package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type participant struct {
	net.Conn
	close chan struct{}
}

type chatRoom struct {
	participants map[string]participant
	lock         sync.Mutex
	chatter      [][]byte
}

func (c *chatRoom) ReadAll() {
	for _, conn := range c.participants {
		go func() {
			reader := bufio.NewReader(conn)
			b, err := reader.ReadBytes('\n')
			if err != nil {
				panic(err) // TODO Work out how to deal with this
				//log.Println(err)
				c.RemovePerson(conn.RemoteAddr().String(), err)
			}
			c.writeAll(conn.RemoteAddr().String(), b)
		}()
	}
}

func (c *chatRoom) welcomeMessage(conn net.Conn) {
	msg := []byte(fmt.Sprintf("Welcome to chatchat, there are %d people here\n", len(c.participants)))
	conn.Write(msg)
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

func (c *chatRoom) AddPerson(p net.Conn, ch chan struct{}) {
	log.Printf("adding chatter %v", p.RemoteAddr().String())
	c.participants[p.RemoteAddr().String()] = participant{p, ch}
}

func (c *chatRoom) RemovePerson(id string, err error) {
	log.Printf("removing chatter %v because %v", id, err)
	ch := c.participants[id].close
	ch <- struct{}{}
	delete(c.participants, id)
}

func main() {
	room := &chatRoom{
		participants: map[string]participant{},
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
			ch := make(chan struct{})
			room.AddPerson(c, ch)
			room.welcomeMessage(c)
			<-ch
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
