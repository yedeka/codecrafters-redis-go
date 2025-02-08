package command

import (
	"strconv"
	"strings"
)

type EchoCommand struct {
	input string
	piggybackFlag bool
	writeCommandFlag bool
}

func (echo EchoCommand) IsReplicaConfigurationAvailabel() (bool, string) {
	return false, ""
}

func (echo EchoCommand) IsReplicationCommand() bool {
	return false
} 

func (echo EchoCommand) IsWriteCommand() bool {
	return echo.writeCommandFlag
}

func (echo EchoCommand) SendPiggyBackResponse() string {
	return noPiggybackResponse
}

func (echo EchoCommand) IsPiggyBackCommand() bool {
	return echo.piggybackFlag 
}

func (echo EchoCommand) Execute() string {
	inputLength := len(echo.input)
	rawResponseList := make([]ParsedResponse, 2)
	rawResponseList[0] = ParsedResponse{
		Responsetype: "LENGTH",
		ResponseData: strconv.Itoa(inputLength),
	}
	rawResponseList[1] = ParsedResponse{
		Responsetype: "DATA",
		ResponseData: echo.input,
	}
	return echo.FormatOutput(rawResponseList)
}

func (echo EchoCommand) FormatOutput(rawResponseList []ParsedResponse) string {
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

func writeResponse(prefix string,
	responseData string,
	responseDelimiter string,
	responseBuffer *strings.Builder) {
	if prefix != "" {
		responseBuffer.WriteString(prefix)
	}
	responseBuffer.WriteString(responseData)
	responseBuffer.WriteString(responseDelimiter)
}
