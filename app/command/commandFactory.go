package command

import (
	"strings"
	"github.com/codecrafters-io/redis-starter-go/app/model"
)

func CommandFactory(inputRequest []string, hostConfig *model.HostConfig) Command {
	argsMap := make(map[string]string)
	
	switch commandName := strings.ToUpper(inputRequest[2]); commandName {
	case "ECHO":
		return EchoCommand{
			input: inputRequest[4],
			piggybackFlag: false,
		}
	case "PING":
		return PingCommand{
			ResponsePrompt: "PONG",
			piggybackFlag: false,
		}
	case "SET":
		keyValue := inputRequest[4]
		if len(inputRequest) > 7 {
			for i := 8; i < len(inputRequest); i += 4 {
				argsMap[strings.ToLower(inputRequest[i])] = inputRequest[i+2]
			}
		}

		return SetCommand{
			key: keyValue,
			value: SetValue{
				key:   keyValue,
				value: inputRequest[6],
			},
			successfulResponse: "+OK",
			args:               argsMap,
			piggybackFlag: false,
		}
	case "GET":
		{
			return GetCommand{
				key:      inputRequest[4],
				erroCode: "-1",
				piggybackFlag: false,
			}
		}
	case "INFO":
		{
			return InfoCommand{
				arguments:  []string{inputRequest[4]},
				hostConfig: hostConfig,
				piggybackFlag: false,
			}
		}
	case "REPLCONF" : 
	{
		return ReplConfCommand{
			arguments:  []string{inputRequest[4]},
			hostConfig: hostConfig,
			piggybackFlag: false,
		}
	}
	case "PSYNC" : 
	{
		//*3 $5 PSYNC $1 ? $2 -1
	
		argsMap["REPL_ID"] = inputRequest[4]
		argsMap["REPL_OFFSET"] = inputRequest[6]
		
		return PSyncCommand{
			arguments:  argsMap,
			hostConfig: hostConfig,
			piggybackFlag: true,
		}
	}
	default:
		return nil
	}
}
