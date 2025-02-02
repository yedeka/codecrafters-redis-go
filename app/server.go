package main

import (
	"flag"
	"github.com/codecrafters-io/redis-starter-go/app/host"
)

var defaultAddress = "0.0.0.0"

func main() {
	var portValue, replicationData string
	flag.StringVar(&portValue, "port", "6379", "Port for Redis server to accept client connections")
	flag.StringVar(&replicationData, "replicaof", "leader", "Replica specification for slave")
	flag.Parse()
	 
	host.CreateHost(portValue, replicationData, defaultAddress)
}

/*func handleConnectionsViaMultiThreading(listener net.Listener) {
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
}*/
