package command

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/model"
)

func CommandFactory(inputRequest []string, hostConfig *model.HostConfig) Command {
	argsMap := make(map[string]string)
	fmt.Printf("%+v\n", inputRequest)
	if len(inputRequest) > 7 {
		for i := 8; i < len(inputRequest); i += 4 {
			argsMap[strings.ToLower(inputRequest[i])] = inputRequest[i+2]
		}
	}
	switch commandName := strings.ToUpper(inputRequest[2]); commandName {
	case "ECHO":
		return EchoCommand{
			input: inputRequest[4],
		}
	case "PING":
		return PingCommand{
			ResponsePrompt: "PONG",
		}
	case "SET":
		keyValue := inputRequest[4]
		return SetCommand{
			key: keyValue,
			value: SetValue{
				key:   keyValue,
				value: inputRequest[6],
			},
			successfulResponse: "+OK",
			args:               argsMap,
		}
	case "GET":
		{
			return GetCommand{
				key:      inputRequest[4],
				erroCode: "-1",
			}
		}
	case "INFO":
		{
			return InfoCommand{
				arguments: []string{inputRequest[4]},
			}
		}

	default:
		return nil
	}
}
