package command

import (
	"strconv"
	"strings"
	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type ReplConfCommand struct {
	hostConfig *model.HostConfig
	arguments  map[string]string
	piggybackFlag bool
	writeCommandFlag bool
	offset int 
}

func (replConf ReplConfCommand) IsReplicationCommand() bool {
	return false
} 


func (replConf ReplConfCommand) IsWriteCommand() bool {
	return replConf.writeCommandFlag
}

func (replConf ReplConfCommand) SendPiggyBackResponse() string {
	return noPiggybackResponse
}

func (replConf ReplConfCommand) IsPiggyBackCommand() bool {
	return replConf.piggybackFlag 
}

func (replConf ReplConfCommand) IsReplicaConfigurationAvailabel() (bool, string) {
	replicationPort, ok := replConf.arguments["listening-port"]; if ok {
		return true, replicationPort 
	}
	return false, ""
} 

func (replConf ReplConfCommand) Execute() string {
	var rawResponseList []ParsedResponse  
	if len(replConf.arguments) == 0 {
		rawResponseList = []ParsedResponse{
			ParsedResponse{
				Responsetype: "SUCCESS",}}
	} else {
		rawResponseList = replConf.handleArguments()
	}
	
	return replConf.FormatOutput(rawResponseList)	
}

func (replConf ReplConfCommand) handleArguments() []ParsedResponse {
	var argumentResponseList []ParsedResponse = []ParsedResponse{} 
	for argKey := range(replConf.arguments) {
		if argKey == GETACK {
			ack_response_tokens := []string{"REPLCONF", "ACK", strconv.Itoa(replConf.offset)}
			argumentResponseList = append(argumentResponseList, prepareParseResponse(ack_response_tokens, true, EMPTY))
		}
	}
	return argumentResponseList
}

func (replConf ReplConfCommand) FormatOutput(rawResponseList []ParsedResponse) string {
	var commandResponse strings.Builder
	for _, rawResponse := range rawResponseList {
		writeResponse(rawResponse.ResponsePrefix,
		rawResponse.ResponseData,
		"",
		&commandResponse)
	}
	return commandResponse.String()
}
