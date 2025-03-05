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

func (client *Follower) GetHostConfig() (*model.HostConfig) {
	return client.hostConfig
}

func (client *Follower) Init() {
	fmt.Println("Initializing Redis Client")
	conn, err := client.connectToServer()
	if nil != err {
		fmt.Printf("Error while connecting to server %s\n", err.Error())
		os.Exit(1)
	}
	client.serverConnection = conn
	err = client.performHandShake()
	if nil == err {
		go HandleReplication(client.serverConnection, client.hostConfig.Offset)
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	} 
}

func (client *Follower) connectToServer() (net.Conn, error) {
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

func (client *Follower) performHandShake() error {
	fmt.Println("Starting Handshake")
	serverResponse := client.sendRequestToServer(pingCommand,  make([]CommandArgument, 0))
	if successfulPingResponse == serverResponse { 
		serverResponse = client.sendRequestToServer(replConfCommand, []CommandArgument {CommandArgument{
			argumentKey: listeningPortConfKey,
			argumentValue: replicationPort,	
		}})
		if successfulResponse == serverResponse { 
			serverResponse = client.sendRequestToServer(replConfCommand, []CommandArgument {
				CommandArgument{
					argumentKey: capacityKey,
					argumentValue: defaultCapacityValue,	
				}})
			if successfulResponse == serverResponse { 
				serverResponse = client.sendRequestToServer(psyncCommand, []CommandArgument {
					CommandArgument{
						argumentKey: emptyKey,
						argumentValue: defaultReplIdValue,	
					}, CommandArgument{
						argumentKey: emptyKey,
						argumentValue: defaultOffsetValue,	
					}})
					if strings.HasPrefix(serverResponse, PSYNCPrefix) {
						return nil
					} else {
						return fmt.Errorf("PSYNC handshake failed with server => %s\n", client.serverConnection.RemoteAddr())		
					}
			} else {
				return fmt.Errorf("REPL-CONF Capacity handshake failed with server => %s\n", client.serverConnection.RemoteAddr())	
			}	 
		} else {
			return fmt.Errorf("REPL-CONF listen-port handshake failed with server => %s\n", client.serverConnection.RemoteAddr())
		}
	}
	return fmt.Errorf("PING handshake failed with server => %s\n", client.serverConnection.RemoteAddr())
}

func (client *Follower) sendRequestToServer(command string, argumentList []CommandArgument) string {
	response := client.sendCommand(command, argumentList)
	client.sendResponseToServer(response)
	serverResponse := client.listenToServer()
	return serverResponse
}

func (client *Follower) sendCommand(command string, argumentList []CommandArgument) string {
	switch commandToSend := command; commandToSend{
	case pingCommand :
		return client.sendSimpleCommand(pingCommand)
	case replConfCommand :
		return client.sendParameterizedCommand(replConfCommand, argumentList)
	case psyncCommand: 
	    return client.sendParameterizedCommand(psyncCommand, argumentList) 	  
	default : 
		return ""
	}
}

func (client *Follower) sendSimpleCommand(simpleCommand string) string {
	outputTokens := []string{}
	outputTokens = append(outputTokens, linePrefix+"1")
	outputTokens = append(outputTokens, lengthPrefix+strconv.Itoa(len(simpleCommand)))
	outputTokens = append(outputTokens, simpleCommand)
	return client.prepOutput(outputTokens)
}

func (client *Follower) sendParameterizedCommand(
	parameterizeCommand string, argumentList []CommandArgument) string {
		outputTokens := []string{}
		outputTokens = append(outputTokens, linePrefix)
		outputTokens = append(outputTokens, 
			lengthPrefix+strconv.Itoa(len(parameterizeCommand)))
		outputTokens = append(outputTokens, parameterizeCommand)
			for _, argument := range(argumentList) {
				if emptyKey != argument.argumentKey {
					outputTokens = append(outputTokens, 
						lengthPrefix+strconv.Itoa(len(argument.argumentKey)))
					outputTokens = append(outputTokens, argument.argumentKey)
				}
				outputTokens = append(outputTokens, 
					lengthPrefix+strconv.Itoa(len(argument.argumentValue)))
				outputTokens = append(outputTokens, argument.argumentValue)

			}
		tokenLength := (len(outputTokens) - 1) / 2
		outputTokens[0] = outputTokens[0] + strconv.Itoa(tokenLength)  
		return client.prepOutput(outputTokens)
}

func (client *Follower) prepOutput(outputTokens []string) string {
	var followerResponse strings.Builder

	for _,token := range(outputTokens) {
		followerResponse.WriteString(token)
		followerResponse.WriteString(delimiter)
	}
	return followerResponse.String()
}

func (client *Follower) sendResponseToServer(responseString string) {
	// Response data to write
	data := []byte(responseString)
	// Write data to the connection
	_, err := client.serverConnection.Write(data)
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}
}

func (client *Follower) listenToServer() string{
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