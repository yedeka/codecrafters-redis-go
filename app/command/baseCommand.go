package command

const terminationSequence = "\r\n"
const lengthPrefix = "$"

type ParsedResponse struct {
	Responsetype string
	ResponseData string
}

type Command interface {
	Execute() string
	FormatOutput([]ParsedResponse) string
}
