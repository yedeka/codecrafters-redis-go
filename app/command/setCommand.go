package command

import "strings"

var keyMap map[string]string

type SetCommand struct {
	key                string
	value              string
	successfulResponse string
}

func (set SetCommand) Execute() string {
	keyMap[set.key] = set.value
	var successfulResponse strings.Builder
	successfulResponse.WriteString(set.successfulResponse)
	successfulResponse.WriteString(terminationSequence)
	return successfulResponse.String()
}

// FormatOutput implements Command.
func (set SetCommand) FormatOutput([]ParsedResponse) string {
	return ("UNUSED")
}
