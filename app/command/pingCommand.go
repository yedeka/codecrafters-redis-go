package command

import (
	"strconv"
	"strings"
)

type PingCommand struct {
	ResponsePrompt string
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
