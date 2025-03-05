package command

import (
	"strconv"
	"strings"
)

const terminationSequence = "\r\n"
const arrayLengthPrefix = "*"
const lengthPrefix = "$"
const rolePrefix = "role:"
const masterRole = "master"
const slaveRole = "slave"
const replictionIdPrefix = "master_replid:"
const replicationOffset = "master_repl_offset:"
const successfulResponse = "+OK"
const fullResyncResponse = "+FULLRESYNC"
const space = " "
const noPiggybackResponse = "Piggyback reponse not implemented"
const GETACK string = "GETACK"
const RESP_ARRAY string = "RESP_ARRAY"
const RESP_STRING string = "RESP_STRING"

const EMPTY string = ""

type ParsedResponse struct {
	ResponsePrefix string
	Responsetype string
	ResponseData string
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

func prepareParseResponse(responseTokens []string, isRespArray bool, responseString string) ParsedResponse {
	respnEncodedResponse := ParsedResponse{}
	if isRespArray {
		respnEncodedResponse.ResponsePrefix = arrayLengthPrefix + strconv.Itoa(len(responseTokens))+terminationSequence
		var responseData strings.Builder
		for _,responseToken := range(responseTokens) {
			responseData.WriteString(lengthPrefix)
			responseData.WriteString(strconv.Itoa(len(responseToken)))
			responseData.WriteString(terminationSequence)
			responseData.WriteString(responseToken)
			responseData.WriteString(terminationSequence)
		}
		respnEncodedResponse.ResponseData = responseData.String()
		respnEncodedResponse.Responsetype = RESP_ARRAY
		return respnEncodedResponse
	}
	respnEncodedResponse.ResponseData = responseString
	respnEncodedResponse.ResponsePrefix = EMPTY
	respnEncodedResponse.Responsetype = RESP_STRING

	return respnEncodedResponse
}