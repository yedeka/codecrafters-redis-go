package command

type Command interface {
	Execute() string
	FormatOutput([]ParsedResponse) string
	IsPiggyBackCommand() bool
	SendPiggyBackResponse() string
	IsWriteCommand() bool
	IsReplicationCommand() bool
	IsReplicaConfigurationAvailabel() (bool, string)
}
