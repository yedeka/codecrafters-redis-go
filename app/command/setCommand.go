package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var keyMap map[string]SetValue = make(map[string]SetValue)

type SetValue struct {
	value string
	key   string
	timer *time.Timer
}

type SetCommand struct {
	key                string
	value              SetValue
	successfulResponse string
	args               map[string]string
}

func (set SetCommand) prepareSetValue() error {
	intDuration, err := strconv.Atoi(set.args["px"])
	if nil != err {
		keyMap[set.key] = set.value
		return errors.New("invalid expiry time specified")
	} else {
		set.value.timer = time.NewTimer(time.Duration(intDuration) * time.Millisecond)
		keyMap[set.key] = set.value
		go func() {
			<-set.value.timer.C
			delete(keyMap, set.key)
		}()
		return nil
	}
}

func (set SetCommand) Execute() string {
	err := set.prepareSetValue()
	if nil != err {
		fmt.Println(err.Error())
	}
	var successfulResponse strings.Builder
	successfulResponse.WriteString(set.successfulResponse)
	successfulResponse.WriteString(terminationSequence)
	return successfulResponse.String()
}

// FormatOutput implements Command.
func (set SetCommand) FormatOutput([]ParsedResponse) string {
	return ("UNUSED")
}
