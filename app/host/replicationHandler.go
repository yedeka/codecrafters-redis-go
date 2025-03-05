package host

import(
	"fmt"
	"net"
	"strings"
	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/util"
)

func HandleReplication(conn net.Conn, offset int) {
	defer conn.Close()
	fmt.Println("Listening to master for commands")	

	for {
		requestData := make([]byte, 1024)
		n, err := conn.Read(requestData)
		if err != nil {
			continue;
		}
		data := requestData[:n]
		serverData := string(data)
		requestBuffer := strings.Split(serverData, "\n")
		inputLength := len(requestBuffer) 

		if strings.HasPrefix(requestBuffer[1], "REDIS") {
			if inputLength > 3 && inputLength > 8  {
				handleReplCommand(requestBuffer[2:], defaultOffset, conn)			
			}
			continue; 
		}  
		handleReplCommand(requestBuffer, defaultOffset, conn)
	}
}

func handleReplCommand(request []string, offset int, connection net.Conn) {
	// Check for first element to be in RESP format since we are not getting that some times. For consistency we add a dummy 
	// element if we do not get a RESP formatted array. 

	receivedCommands := command.ReplicationCommandFactory(request, offset)
	for _,command := range(receivedCommands) {
		commandResponse := command.ReplCommand.Execute()
		if command.IsResponseAvailable {
			err := util.SendResponseOverConnection(connection, commandResponse)
			if nil != err {
				fmt.Printf("Error handling replication command => %s", err)
			}
		}   
	}  
}
