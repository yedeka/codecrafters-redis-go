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
			writeCommandFlag: false,
		}
	case "PING":
		return PingCommand{
			ResponsePrompt: "PONG",
			piggybackFlag: false,
			writeCommandFlag: false,
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
			writeCommandFlag: true,
		}
	case "GET":
		{
			return GetCommand{
				key:      inputRequest[4],
				erroCode: "-1",
				piggybackFlag: false,
				writeCommandFlag: false,
			}
		}
	case "INFO":
		{
			return InfoCommand{
				arguments:  []string{inputRequest[4]},
				hostConfig: hostConfig,
				piggybackFlag: false,
				writeCommandFlag: false,
			}
		}
	case "REPLCONF" : 
	{
		argsMap[inputRequest[4]] = inputRequest[6] 
		
		return ReplConfCommand{
			arguments:  argsMap,
			hostConfig: hostConfig,
			piggybackFlag: false,
			writeCommandFlag: false,
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
			writeCommandFlag: false,
		}
	}
	default:
		return nil
	}
}

func ReplicationCommandFactory(request []string, followerOffset int) []ReplicationCommand {
	commandArray := []ReplicationCommand{}
	switch commandName := strings.TrimSpace(strings.ToUpper(request[2])); commandName {
	case "SET": 
		for count:=0; count<len(request); count+=7{
			if len(strings.TrimSpace(request[count])) > 0 {
				commandArray = append(commandArray, ReplicationCommand{
					ReplCommand : SetCommand{
						key: request[count+4],
						value: SetValue{
							key:   request[count+4],
							value: request[count+6],
						},
						successfulResponse: "+OK",
						args:               make(map[string]string),
						piggybackFlag: false,
						writeCommandFlag: false,
					},
					IsResponseAvailable: false, 
				})
			}	
		}
		return commandArray

	case "REPLCONF" : 
	{
		argsMap := make(map[string]string)
		argsMap[strings.TrimSpace(request[4])] = strings.TrimSpace(request[6])
		commandArray = append(commandArray, ReplicationCommand{
			ReplCommand : ReplConfCommand{
				arguments:  argsMap,
				piggybackFlag: false,
				writeCommandFlag: false,
				offset: followerOffset,
			}, 
			IsResponseAvailable: true,
		})
		return commandArray
	}	
	default:
		return nil
	}	  
}
