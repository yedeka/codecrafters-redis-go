package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type InfoCommand struct {
	hostConfig *model.HostConfig
	arguments  []string
}

func handleArgs(info InfoCommand) []ParsedResponse {
	responseList := []ParsedResponse{}
	for _, arg := range info.arguments {
		if strings.ToLower(arg) == "replication" {
			var responseString string = ""
			if info.hostConfig.IsMaster {
				return handleMaster(defaultReplId, defaultOffset)
			} else {
				responseString = "role:slave"
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
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "LENGTH",
		ResponseData: strconv.Itoa(len(replicationId)),
	})
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "REPL_ID",
		ResponseData: replicationId,
	})
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "LENGTH",
		ResponseData: strconv.Itoa(len(replicationId)),
	})
	masterInfoResponse = append(masterInfoResponse, ParsedResponse{
		Responsetype: "REPL_OFFSET",
		ResponseData: replicationId,
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
		} else if rawResponse.Responsetype == "REPL_ID" {
			writeResponse(replictionIdPrefix,
				rawResponse.ResponseData,
				terminationSequence,
				&commandResponse)
		} else if rawResponse.Responsetype == "REPL_OFFSET" {
			writeResponse(replicationOffset,
				rawResponse.ResponseData,
				terminationSequence,
				&commandResponse)
		} else if rawResponse.Responsetype == "DATA" {
			writeResponse("",
				rawResponse.ResponseData,
				terminationSequence,
				&commandResponse)
		}
	}
	fmt.Printf("CommandREsponse => %s", commandResponse.String())
	return commandResponse.String()
}
