package command

import (
	"strconv"
	"strings"
)

type PingCommand struct {
	ResponsePrompt string
	piggybackFlag bool
	writeCommandFlag bool
}

func (ping PingCommand) IsReplicationCommand() bool {
	return false
} 

func (ping PingCommand) SendPiggyBackResponse() string {
	return noPiggybackResponse
}

func (ping PingCommand) IsWriteCommand() bool {
	return ping.writeCommandFlag
}

func (ping PingCommand) IsPiggyBackCommand() bool {
	return ping.piggybackFlag 
}
// FormatOutput implements Command.
func (ping PingCommand) FormatOutput([]ParsedResponse) string {
	return ("UNUSED")
}

func (ping PingCommand) Execute() string {
	var pingResponse strings.Builder
	lengthString := strconv.Itoa(len(ping.ResponsePrompt))
	// Writing Ping response here without a formatter since it is a simple static response
	pingResponse.WriteString(lengthPrefix)
	pingResponse.WriteString(lengthString)
	pingResponse.WriteString(terminationSequence)
	pingResponse.WriteString(ping.ResponsePrompt)
	pingResponse.WriteString(terminationSequence)

	return pingResponse.String()
}
