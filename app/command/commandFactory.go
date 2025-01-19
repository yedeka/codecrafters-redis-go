package command

import "strings"

func CommandFactory(inputRequest []string) Command {
	argsMap := make(map[string]string)
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
		return SetCommand{
			key: SetKey{
				key: inputRequest[4],
			},
			value:              inputRequest[6],
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

	default:
		return nil
	}
}
