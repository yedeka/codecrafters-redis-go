package command

import (
	"fmt"
	"strconv"
	"strings"
)

type GetCommand struct {
	key      string
	erroCode string
}

func (get GetCommand) Execute() string {
	var rawResponseList []ParsedResponse
	value, ok := keyMap[get.key]
	fmt.Printf("%+v\n", keyMap)
	if !ok {
		rawResponseList := append(rawResponseList, ParsedResponse{
			Responsetype: "FAILURE",
			ResponseData: get.erroCode,
		})
		return get.FormatOutput(rawResponseList)
	} else {
		rawResponseList = append(rawResponseList, ParsedResponse{
			Responsetype: "LENGTH",
			ResponseData: strconv.Itoa(len(value.value)),
		})
		rawResponseList = append(rawResponseList, ParsedResponse{
			Responsetype: "DATA",
			ResponseData: value.value,
		})
		return get.FormatOutput(rawResponseList)
	}

}

// FormatOutput implements Command.
func (get GetCommand) FormatOutput(rawResponseList []ParsedResponse) string {
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
		} else if rawResponse.Responsetype == "FAILURE" {
			writeResponse(lengthPrefix,
				rawResponse.ResponseData,
				terminationSequence,
				&commandResponse)
		}
	}

	return commandResponse.String()
}
