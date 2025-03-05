package util 

import (
	"net" 
	"fmt"
)

func SendResponseOverConnection(connection net.Conn, serverResponse string) error{
	// Write back to connection
	_, err := connection.Write([]byte(serverResponse))
	if err != nil {
		fmt.Println("Could not write back to channel")
		return err
	}
	return nil
}	