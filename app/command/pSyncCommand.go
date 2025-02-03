package command

import (
	"io/ioutil"
	"strconv"
	"strings"
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/model"
)

type PSyncCommand struct {
	hostConfig *model.HostConfig
	arguments  map[string]string
	piggybackFlag bool
}

func (psync PSyncCommand) SendPiggyBackResponse() string {
	rdbContent, err := ioutil.ReadFile("empty.rdb")
	if err != nil {
		log.Fatal(err)
	}
	// Print the byte slice
	rdbContentString := string(rdbContent)
	var emptyRdbResponse strings.Builder
	emptyRdbResponse.WriteString(lengthPrefix)
	emptyRdbResponse.WriteString(strconv.Itoa(len(rdbContentString)))
	emptyRdbResponse.WriteString(terminationSequence)
	emptyRdbResponse.WriteString(rdbContentString)
	return emptyRdbResponse.String()
}

func (psync PSyncCommand) IsPiggyBackCommand() bool {
	return psync.piggybackFlag 
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
