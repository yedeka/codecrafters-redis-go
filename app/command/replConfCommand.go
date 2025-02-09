package command

import (
	"strings"
	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type ReplConfCommand struct {
	hostConfig *model.HostConfig
	arguments  map[string]string
	piggybackFlag bool
	writeCommandFlag bool
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
