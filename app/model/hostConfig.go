package model

type HostConfig struct {
	IsMaster    bool
	MasterProps MasterConfig
	ReplId string
	Offset int
}
