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
	piggybackFlag bool
	writeCommandFlag bool
}

func (set SetCommand) IsWriteCommand() bool {
	return set.writeCommandFlag
}

func (set SetCommand) SendPiggyBackResponse() string {
	return noPiggybackResponse
}

func (set SetCommand) IsPiggyBackCommand() bool {
	return set.piggybackFlag 
}

func (set SetCommand) prepareSetValue() error {
	if pxValue, ok := set.args["px"]; ok { 
		intDuration, err := strconv.Atoi(pxValue)
		if nil != err { 
			return errors.New("invalid expiry time specified")
		}
		set.value.timer = time.NewTimer(time.Duration(intDuration) * time.Millisecond)
		keyMap[set.key] = set.value
		go func() {
			<-set.value.timer.C
			delete(keyMap, set.key)
		}()
		return nil
	}
    keyMap[set.key] = set.value
	return nil 
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
