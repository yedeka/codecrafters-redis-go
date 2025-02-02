package host

import (
	"net"
	"fmt"
	"os"
	"strings"
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type Master struct {
	hostIpAddress string 
	listeningPort string
	hostConfig *model.HostConfig
}

func (master Master) GetHostConfig() *model.HostConfig {
	return master.hostConfig
}

func (master Master) Init() {
	var listeningAddress strings.Builder
	listeningAddress.WriteString(master.hostIpAddress)
	listeningAddress.WriteString(serverDelimiter)
	listeningAddress.WriteString(master.listeningPort)

	master.listenForConnections(listeningAddress.String())
}

func (master Master) listenForConnections(listeningAddress string){
	l, err := net.Listen("tcp", listeningAddress)
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", master.listeningPort)
		os.Exit(1)
	}
	defer l.Close()
	master.handleConnectionsViaEventLoop(l)
}

// handleConnectionsViaEventLoop tries to implement an Event Loop like structure to handle the connection requests by trying to keep
// main execution thread free.

func (master Master) handleConnectionsViaEventLoop(listener net.Listener) {
	// Create a channel to listen for and populate a new request
	connChannel := make(chan net.Conn)
	go master.acceptConnections(listener, connChannel)
	// The only waiting call here is the call that now waits for connections to be populated
	for {
		// Capture the connection object coming from channel here in a variable.
		// This will wait until the connection object gets populated in `acceptConnections` function.
		// We need a for loop here since we have to keep on listening for the channel contineously to respond to incoming requests.
		conn := <-connChannel
		go master.handleCons(conn)
	}
}

// acceptConnections accepts connections on the specified port and sends the new connection to a channel for further processing.  
func (master Master) acceptConnections(listener net.Listener, acceptChan chan net.Conn) {
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


// handleCons handles the connection by parsing the incoming command and responding appropriately
func (master Master) handleCons(conn net.Conn) {
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
		requestedCommand := command.CommandFactory(requestBuffer, master.hostConfig)
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
