package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/command"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var defaultAddress = "0.0.0.0:"

func main() {
	var portFlag, replicaFlag string
	flag.StringVar(&portFlag, "port", "6379", "Port for Redis server to accept client connections")
	flag.StringVar(&replicaFlag, "replicaof", "localhost 6380", "Replica specification for slave")
	flag.Parse()
	var listeningAddress strings.Builder
	var replicaList = strings.Split(replicaFlag, " ")
	fmt.Printf("%+v\n", replicaList)
	listeningAddress.WriteString(defaultAddress)
	listeningAddress.WriteString(portFlag)
	l, err := net.Listen("tcp", listeningAddress.String())

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	//handleConnectionsViaMultiThreading(l)
	handleConnectionsViaEventLoop(l)
}

// handleConnectionsViaEventLoop tries to implement an Event Loop like structure to handle the connection requests by trying to keep
// main execution thread free.
func handleConnectionsViaEventLoop(listener net.Listener) {
	// Create a channel to listen for and populate a new request
	connChannel := make(chan net.Conn)
	go acceptConnections(listener, connChannel)
	// The only waiting call here is the call that now waits for connections to be populated
	for {
		// Capture the connection object coming from channel here in a variable.
		// This will wait until the connection object gets populated in `acceptConnections` function.
		// We need a for loop here since we have to keep on listening for the channel contineously to respond to incoming requests.
		conn := <-connChannel
		go handleConns(conn)
	}
}

func acceptConnections(listener net.Listener, acceptChan chan net.Conn) {
	// This blocking behavior is necessary to keep on listening on given port for new connection requests from Clients
	// How main thread will be unblocked here is by invoking the acceptConnections function as a go subroutine.
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			conn.Close()
			continue
		}
		acceptChan <- conn
	}
}

func handleConns(conn net.Conn) {
	defer conn.Close()
	var requestBuffer []string
	for {
		requestData := make([]byte, 1024)
		n, err := conn.Read(requestData)
		if err != nil {
			continue
		}
		data := requestData[:n]
		requestBuffer = strings.Split(string(data), "\r\n")
		requestedCommand := command.CommandFactory(requestBuffer)
		if nil == requestedCommand {
			fmt.Println("Unsupported command passed")
			os.Exit(1)
		}
		// Write back to connection
		_, err = conn.Write([]byte(requestedCommand.Execute()))
		if err != nil {
			fmt.Println("Could not write back to channel")
			continue
		}
	}
}

func handleConnectionsViaMultiThreading(listener net.Listener) {
	for {
		// Accepting connections in an infinite for loop so that we can accept connections from multiple clients.
		// DANGER - Blocking call here we will block until a new client requests connection.
		conn, err := listener.Accept()
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
		// DANGER: Blocking call here waiting for a client to request for something on a connection.
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			return
		}
		// Write back to connection
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			return
		}
	}
}
