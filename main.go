package main

import (
	"fmt"
	"log"
	"net"
)

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
		log.Println("Receivedd message: ", []byte(msg))

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
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}
