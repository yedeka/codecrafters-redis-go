package command

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var keyMap map[SetKey]string = make(map[SetKey]string)

type SetKey struct {
	key   string
	timer *time.Timer
}

type SetCommand struct {
	key                SetKey
	value              string
	successfulResponse string
	args               map[string]string
}

func (set SetCommand) prepareSetKey() error {
	intDuration, err := strconv.Atoi(set.args["px"])
	if nil != err {
		errors.New("Invalid expiry time specified")
	}
	set.key.timer = time.NewTimer(time.Duration(intDuration) * time.Millisecond)
	go func() {
		<-set.key.timer.C
		delete(keyMap, set.key)
	}()
	return nil
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
