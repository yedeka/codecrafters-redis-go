package command

import (
	"strconv"
	"strings"
)

type InfoCommand struct {
	arguments []string
}

func handleArgs(args []string) []ParsedResponse {
	responseList := []ParsedResponse{}
	for _, arg := range args {
		if strings.ToLower(arg) == "replication" {
			responseString := "role:master"
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
	return info.FormatOutput(handleArgs(info.arguments))
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
