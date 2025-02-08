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
	ackPendingReplica map[net.Addr]string
	replicaConnections []net.Conn
}

func (master *Master) GetHostConfig() *model.HostConfig {
	return master.hostConfig
}

func (master *Master) Init() {
	master.replicaConnections = []net.Conn{}
	master.ackPendingReplica = map[net.Addr]string{}
	var listeningAddress strings.Builder
	listeningAddress.WriteString(master.hostIpAddress)
	listeningAddress.WriteString(serverDelimiter)
	listeningAddress.WriteString(master.listeningPort)

	master.listenForConnections(listeningAddress.String())
}

func (master *Master) listenForConnections(listeningAddress string){
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

func (master *Master) handleConnectionsViaEventLoop(listener net.Listener) {
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
func (master *Master) acceptConnections(listener net.Listener, acceptChan chan net.Conn) {
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
func (master *Master) handleCons(conn net.Conn) {
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
		respnEncodedResponse := requestedCommand.Execute()
		err = writeDataToConnection(conn, respnEncodedResponse)
		if nil != err {
			fmt.Println("Could not write data to connection")
		} else {
			if requestedCommand.IsWriteCommand() {
				fmt.Println("got write command sending over data to replicas")
				// Send data over to all the replicas for writing		
				for _, replicaConn := range(master.replicaConnections) {
					fmt.Printf("Sending data to replication connection %s\n", respnEncodedResponse)
					writeDataToConnection(replicaConn, respnEncodedResponse)
				}					
			}
			if requestedCommand.IsPiggyBackCommand() {
				err = writeDataToConnection(conn, requestedCommand.SendPiggyBackResponse())
				if nil != err {
					fmt.Printf("Error while sending Piggyback response to client %s",err.Error())
				}
				err = createReplicationConnection(master, conn)
				if(nil != err) {
					fmt.Println(err.Error())
				}
			}
			configAvailable, replicationPort := requestedCommand.IsReplicaConfigurationAvailabel()
			if configAvailable {
				master.ackPendingReplica[conn.RemoteAddr()] = replicationPort
			}
		}
	}
}

func createReplicationConnection(master *Master, conn net.Conn) error {
	_, ok := master.ackPendingReplica[conn.RemoteAddr()]
	if !ok {
		return fmt.Errorf("error : Sending RDB file to a connection outside Handshake")
	}
	master.replicaConnections = append(master.replicaConnections, conn)
	delete(master.ackPendingReplica, conn.RemoteAddr())
	return nil
}

func writeDataToConnection(connection net.Conn, serverResponse string) error{
	// Write back to connection
	_, err := connection.Write([]byte(serverResponse))
	if err != nil {
		fmt.Println("Could not write back to channel")
		return err
	}
	return nil
}
