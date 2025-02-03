package command

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type InfoCommand struct {
	hostConfig *model.HostConfig
	arguments  []string
	piggybackFlag bool
	writeCommandFlag bool
}

func (info InfoCommand) IsReplicationCommand() bool {
	return false
} 

func (info InfoCommand) IsWriteCommand() bool {
	return info.writeCommandFlag
}

func (info InfoCommand) IsPiggyBackCommand() bool {
	return info.piggybackFlag 
}

func (info InfoCommand) SendPiggyBackResponse() string {
	return noPiggybackResponse
}

func handleArgs(info InfoCommand) []ParsedResponse {
	responseList := []ParsedResponse{}
	for _, arg := range info.arguments {
		if strings.ToLower(arg) == "replication" {
			var responseString string = ""
			if info.hostConfig.IsMaster {
				return handleMaster(info.hostConfig.ReplId, info.hostConfig.Offset)
			} else {
				responseString = rolePrefix + slaveRole
			}
			responseList = append(responseList, ParsedResponse{
				Responsetype: "LENGTH",
				ResponseData: strconv.Itoa(len(responseString)),
			})
			responseList = append(responseList, ParsedResponse{
				Responsetype: "DATA",
				ResponseData: responseString,
			})
		}
	}
	return responseList
}

func handleMaster(replicationId string, offset int) []ParsedResponse {
	//8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb
	masterInfoResponse := []ParsedResponse{}
	hostRoleEntry := rolePrefix + masterRole
	replicationIdEntry := replictionIdPrefix + replicationId
	replicationOffset := replicationOffset + strconv.Itoa(offset)

	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "LENGTH",
		ResponseData: strconv.Itoa(len(hostRoleEntry+replicationIdEntry+replicationOffset) + 4),
	})
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "HOST_ROLE",
		ResponseData: hostRoleEntry,
	})
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "REPL_ID",
		ResponseData: replicationIdEntry,
	})
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "REPL_OFFSET",
		ResponseData: replicationOffset,
	})
	return masterInfoResponse
}

func (info InfoCommand) Execute() string {
	return info.FormatOutput(handleArgs(info))
}

func (info InfoCommand) FormatOutput(rawResponseList []ParsedResponse) string {
	var commandResponse strings.Builder
	for _, rawResponse := range rawResponseList {
		if rawResponse.Responsetype == "LENGTH" {
			writeResponse(lengthPrefix,
				rawResponse.ResponseData,
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
