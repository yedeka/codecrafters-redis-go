package host

import(
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
	response := client.sendCommand(pingCommand)
	client.sendResponseToServer(response)
}

func (client Follower) sendCommand(command string) string {
	//*1\r\n$4\r\nPING\r\n
	outputTokens := []string{}
	outputTokens = append(outputTokens, linePrefix+"1")
	outputTokens = append(outputTokens, lengthPrefix+strconv.Itoa(len(command)))
	outputTokens = append(outputTokens, command)
	fmt.Printf("%v\n", outputTokens)
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