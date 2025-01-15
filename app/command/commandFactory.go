package command

import "strings"

func CommandFactory(inputRequest []string) Command {
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
		{
			return SetCommand{
				key:                inputRequest[4],
				value:              inputRequest[6],
				successfulResponse: "+OK",
			}
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
