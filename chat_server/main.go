package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	conns   []net.Conn
	connCh  = make(chan net.Conn)
	closeCh = make(chan net.Conn)
	msgCh   = make(chan string)
)

func main() {
	fmt.Println("Start chat server")
	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		fmt.Println("Start listening connection")
		for {
			conn, err := server.Accept()
			fmt.Println("Clent connected")

			if err != nil {
				log.Fatal(err)
			}
			connCh <- conn
			conns = append(conns, conn)
		}
	}()

	for {
		select {
		case conn := <-connCh:
			go onMessage(conn)
		case msg := <-msgCh:
			fmt.Println(strings.Trim(msg, "\n"))

		case conn := <-closeCh:
			fmt.Println("Client Exit")
			removeConn(conn)
		}

	}
}

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msgCh <- msg
		pulishMsg(conn, msg)

	}
	closeCh <- conn
}

func removeConn(conn net.Conn) {
	var i int
	for i = range conns {
		if conn == conns[i] {
			break
		}
	}
	conns = append(conns[i:], conns[:i+1]...)
}

func pulishMsg(conn net.Conn, msg string) {
	for i := range conns {
		if conn != conns[i] {
			conns[i].Write([]byte(msg))
		}
	}
}
