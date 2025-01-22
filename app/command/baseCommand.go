package command

const terminationSequence = "\r\n"
const lengthPrefix = "$"
const replictionIdPrefix = "master_replid"
const replicationOffset = "master_repl_offset"
const defaultReplId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
const defaultOffset = 0

type ParsedResponse struct {
	Responsetype string
	ResponseData string
}

type Command interface {
	Execute() string
	FormatOutput([]ParsedResponse) string
}
