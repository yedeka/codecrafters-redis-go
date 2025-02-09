package host

import(
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/command"
)

func HandleReplication(conn net.Conn) {
	for {
		requestData := make([]byte, 1024)
		n, err := conn.Read(requestData)
		if err != nil {
			fmt.Printf("Error => %s",err.Error())
			continue;
		}
		data := requestData[:n]
		requestBuffer := strings.Split(string(data), "\r\n")
		// Check for RDB file contents
		if strings.HasPrefix(requestBuffer[1], "REDIS") {
			continue;
		}
		go handleReplCommand(requestBuffer)
	}
}

func handleReplCommand(request []string) {
	receivedCommands := command.ReplicationCommandFactory(request)
	fmt.Printf("%+v\n", receivedCommands)
	for _,command := range(receivedCommands) {
		commandResponse := command.ReplCommand.Execute()
		fmt.Println(commandResponse)  
	}  
}
