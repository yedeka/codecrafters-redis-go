package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	handleConnectionsViaMultiThreading(l)
}

func handleConnectionsViaMultiThreading(listener net.Listener) {
	for {
		// Accepting connections in an infinite for loop so that we can accept connections from multiple clients.
		// DANGER - Blocking call here we will block until a new client requests connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			conn.Close()
			os.Exit(1)
		}
		// Handle connections asynchroneously so that multiple connections can be responded accordingly.
		go handleCommand(conn)
	}
}

func handleCommand(conn net.Conn) {
	defer conn.Close()
	for {
		// Read data from the connection
		buf := make([]byte, 1024)
		im, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Printf("Incoming message %d\n", im)
		// Write back to connection
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			return
		}
	}
}
