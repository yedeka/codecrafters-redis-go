package command

const terminationSequence = "\r\n"
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

type ParsedResponse struct {
	Responsetype string
	ResponseData string
}

type Command interface {
	Execute() string
	FormatOutput([]ParsedResponse) string
	IsPiggyBackCommand() bool
	SendPiggyBackResponse() string
	IsWriteCommand() bool
	IsReplicationCommand() bool
	IsReplicaConfigurationAvailabel() (bool, string)
}
