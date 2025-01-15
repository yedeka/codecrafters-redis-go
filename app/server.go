package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

const terminationSequence = "\r\n"
const lengthPrefix = "$"

type ParsedResponse struct {
	Responsetype string
	ResponseData string
}

type Command interface {
	Execute() string
	FormatOutput([]ParsedResponse) string
}

type EchoCommand struct {
	input string
}

func (echo EchoCommand) Execute() string {
	inputLength := len(echo.input)
	fmt.Println("input ", echo.input)
	rawResponseList := make([]ParsedResponse, 2)
	rawResponseList[0] = ParsedResponse{
		Responsetype: "LENGTH",
		ResponseData: strconv.Itoa(inputLength),
	}
	rawResponseList[1] = ParsedResponse{
		Responsetype: "DATA",
		ResponseData: echo.input,
	}
	return echo.FormatOutput(rawResponseList)
}

func (echo EchoCommand) FormatOutput(rawResponseList []ParsedResponse) string {
	var commandResponse strings.Builder
	for _, rawResponse := range rawResponseList {
		if rawResponse.Responsetype == "LENGTH" {
			writeResponse(lengthPrefix,
				rawResponse.ResponseData,
				terminationSequence,
				commandResponse)
		} else if rawResponse.Responsetype == "DATA" {
			writeResponse("",
				rawResponse.ResponseData,
				terminationSequence,
				commandResponse)
		}
	}
	return commandResponse.String()
}

func writeResponse(prefix string,
	responseData string,
	responseDelimiter string,
	responseBuffer strings.Builder) {
	if prefix != "" {
		responseBuffer.WriteString(prefix)
	}
	responseBuffer.WriteString(responseData)
	responseBuffer.WriteString(responseDelimiter)
	fmt.Println(responseBuffer.String() + "\n")
}

type PingCommand struct {
	ResponsePrompt string
}

func (ping PingCommand) Execute() string {
	var pingResponse strings.Builder
	lengthString := strconv.Itoa(len(ping.ResponsePrompt))
	// Writing Ping response here without a formatter since it is a simple static response
	pingResponse.WriteString(lengthPrefix)
	pingResponse.WriteString(lengthString)
	pingResponse.WriteString(terminationSequence)
	pingResponse.WriteString(ping.ResponsePrompt)
	pingResponse.WriteString(terminationSequence)

	return pingResponse.String()
}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
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

func handleRequest(inputRequest []string) string {
	switch command := strings.ToUpper(inputRequest[2]); command {
	case "ECHO":
		return EchoCommand{
			input: inputRequest[4],
		}.Execute()
	case "PING":
		return PingCommand{
			ResponsePrompt: "PONG",
		}.Execute()
	default:
		return "Non Standard Command"
	}
}

func handleConns(conn net.Conn) {
	defer conn.Close()
	//reader := bufio.NewReader(conn)
	var requestBuffer []string
	for {
		requestData := make([]byte, 1024)
		n, err := conn.Read(requestData)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		data := requestData[:n]
		requestBuffer = strings.Split(string(data), "\r\n")
		// Write back to connection
		_, err = conn.Write([]byte(handleRequest(requestBuffer)))
		if err != nil {
			fmt.Println("Could not write back to channel")
			continue
		}
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
