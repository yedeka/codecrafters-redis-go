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
				key:   inputRequest[4],
				value: inputRequest[6],
			}
		}
	default:
		return nil
	}
}
