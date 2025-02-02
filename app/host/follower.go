package host

import(
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	
	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type Follower struct {
	hostConfig *model.HostConfig
	serverConnection net.Conn
}

type CommandArgument struct {
	argumentKey string
	argumentValue string
}

func (client Follower) GetHostConfig() (*model.HostConfig) {
	return client.hostConfig
}

func (client Follower) Init() {
	fmt.Println("Initializing Redis Client")
	// client.serverConnection = 
	conn, err := client.connectToServer()
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	client.serverConnection = conn
	client.performHandShake()
	defer client.serverConnection.Close()
}

func (client Follower) connectToServer() (net.Conn, error) {
	var connectionString strings.Builder
	connectionString.WriteString(client.hostConfig.MasterProps.Host)
	connectionString.WriteString(serverDelimiter)
	connectionString.WriteString(strconv.Itoa(client.hostConfig.MasterProps.Port))
	connectionURL := connectionString.String()
	// Connect to server
	fmt.Printf("Connecting to server %s\n", connectionURL)
	conn, err := net.Dial("tcp", connectionURL)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to server %s", connectionURL)
	}
	return conn, nil
}

func (client Follower) performHandShake() {
	fmt.Println("Starting Handshake")
	serverResponse := client.sendRequestToServer(pingCommand, CommandArgument{})
	if successfulPingResponse == serverResponse { 
		fmt.Println("Successful PING handshake")
		serverResponse = client.sendRequestToServer(replConfCommand, CommandArgument{
			argumentKey: listeningPortConfKey,
			argumentValue: replicationPort,	
		})
		if successfulResponse == serverResponse { 
			fmt.Println("REPLCONF listen-port complete")
			serverResponse = client.sendRequestToServer(replConfCommand, CommandArgument{
				argumentKey: capacityKey,
				argumentValue: defaultCapacityValue,	
			})
			if successfulResponse == serverResponse { 
				fmt.Println("Succeful completion of REPLCONF handshake")
			}	
		}
	}
}

func (client Follower) sendRequestToServer(command string, argument CommandArgument) string {
	response := client.sendCommand(command, argument)
	client.sendResponseToServer(response)
	serverResponse := client.listenToServer()
	return serverResponse
}

func (client Follower) sendCommand(command string, argument CommandArgument) string {
	switch commandToSend := command; commandToSend{
	case pingCommand :
		return client.sendSimpleCommand(pingCommand)
	case replConfCommand :
		return client.sendParameterizedCommand(replConfCommand, argument) 	  
	default : 
		return ""
	}
}

func (client Follower) sendSimpleCommand(simpleCommand string) string {
	outputTokens := []string{}
	outputTokens = append(outputTokens, linePrefix+"1")
	outputTokens = append(outputTokens, lengthPrefix+strconv.Itoa(len(simpleCommand)))
	outputTokens = append(outputTokens, simpleCommand)
	return client.prepOutput(outputTokens)
}

func (client Follower) sendParameterizedCommand(
	parameterizeCommand string, argument CommandArgument) string{
		outputTokens := []string{}
		outputTokens = append(outputTokens, linePrefix+"3")
		outputTokens = append(outputTokens, 
			lengthPrefix+strconv.Itoa(len(parameterizeCommand)))
		outputTokens = append(outputTokens, parameterizeCommand)
		outputTokens = append(outputTokens, 
			lengthPrefix+strconv.Itoa(len(argument.argumentKey)))
		outputTokens = append(outputTokens, argument.argumentKey)
		outputTokens = append(outputTokens, 
			lengthPrefix+strconv.Itoa(len(argument.argumentValue)))
		outputTokens = append(outputTokens, argument.argumentValue)
		return client.prepOutput(outputTokens)
}

func (client Follower) prepOutput(outputTokens []string) string {
	var followerResponse strings.Builder

	for _,token := range(outputTokens) {
		followerResponse.WriteString(token)
		followerResponse.WriteString(delimiter)
	}
	return followerResponse.String()
}

func (client Follower) sendResponseToServer(responseString string) {
	// Response data to write
	data := []byte(responseString)
	// Write data to the connection
	_, err := client.serverConnection.Write(data)
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
}

func (client Follower) listenToServer() string{
	// Read response from the server
	reader := bufio.NewReader(client.serverConnection)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		return message
	}
	return " "
}