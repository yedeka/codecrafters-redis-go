package command

import (
	"fmt"
	"strings"
)

var keyMap map[string]string = make(map[string]string)

type SetCommand struct {
	key                string
	value              string
	successfulResponse string
	args               map[string]string
}

func (set SetCommand) Execute() string {
	fmt.Printf("argsMap => %+v", set.args)
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
