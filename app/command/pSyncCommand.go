package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type PSyncCommand struct {
	hostConfig *model.HostConfig
	arguments  map[string]string
}


func (psync PSyncCommand) Execute() string {
	rawResponseList := []ParsedResponse{
		ParsedResponse{
			Responsetype: "FULLRESYNC",}}
	return psync.FormatOutput(rawResponseList)	
}

func (psync PSyncCommand) FormatOutput(rawResponseList []ParsedResponse) string {
	var commandResponse strings.Builder
	for _, rawResponse := range rawResponseList {
		if rawResponse.Responsetype == "FULLRESYNC" {
			var resyncResponse strings.Builder
			resyncResponse.WriteString(fullResyncResponse)
			resyncResponse.WriteString(space)
			resyncResponse.WriteString(psync.hostConfig.ReplId)
			resyncResponse.WriteString(space)
			resyncResponse.WriteString(strconv.Itoa(psync.hostConfig.Offset))
			writeResponse("",
			resyncResponse.String(),
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
