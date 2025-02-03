package command

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type ReplConfCommand struct {
	hostConfig *model.HostConfig
	arguments  []string
	piggybackFlag bool
}

func (replConf ReplConfCommand) SendPiggyBackResponse() string {
	return noPiggybackResponse
}

func (replConf ReplConfCommand) IsPiggyBackCommand() bool {
	return replConf.piggybackFlag 
}

func (replConf ReplConfCommand) Execute() string {
	rawResponseList := []ParsedResponse{
		ParsedResponse{
			Responsetype: "SUCCESS",}}
	return replConf.FormatOutput(rawResponseList)	
}

func (replConf ReplConfCommand) FormatOutput(rawResponseList []ParsedResponse) string {
	var commandResponse strings.Builder
	for _, rawResponse := range rawResponseList {
		if rawResponse.Responsetype == "SUCCESS" {
			writeResponse("",
				successfulResponse,
				terminationSequence,
				&commandResponse)
		} else {
			writeResponse("",
				rawResponse.ResponseData,
				terminationSequence,
				&commandResponse)
		}
	}
	return commandResponse.String()
}
