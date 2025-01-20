package command

import (
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
				responseString = "role:master"
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
		} else if rawResponse.Responsetype == "DATA" {
			writeResponse("",
				rawResponse.ResponseData,
				terminationSequence,
				&commandResponse)
		}
	}

	return commandResponse.String()
}
