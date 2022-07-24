package main

import (
	"fmt"
	"log"
	"net"
)

type UserJoinedEvent struct {
	user *User
}

type User struct {
	name    string
	session *Session
}

type Session struct {
	conn *net.Conn
}

type World struct {
	users []*User
}

func handleConnection(conn net.Conn) error {
	log.Println("I got a connection")
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println(buf)
			break
		}
		msg := buf[0 : n-2]
		log.Println("Recieved message: ", []byte(msg))

		resp := fmt.Sprintf("You said, \"%s\"\r\n", msg)
		n, err = conn.Write([]byte(resp))
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			break
		}
	}
	return nil
}

func startServer() error {
	log.Println("Starting server")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConnection(conn); err != nil {
				log.Panicln("Error handling connection", err)
				return
			}
		}()
	}
}

func main() {
	w := &World{}

	ch := make(chan interface{})

	go func() {
		err := startServer(ch)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
