package host

import("github.com/codecrafters-io/redis-starter-go/app/model")
// Common constants used for sending RESP array responses
const linePrefix = "*"
const delimiter = "\r\n"
const lengthPrefix = "$"
const serverDelimiter = ":"
// Constants for command names
const pingCommand = "PING"
const replConfCommand = "REPLCONF"
const psyncCommand = "PSYNC"
// Constants for command arguments
const listeningPortConfKey = "listening-port"
const replicationPort = "6380"
const capacityKey = "capa"
const defaultCapacityValue = "psync2"
const emptyKey = "EMPTY"
const defaultReplIdValue = "?"
const defaultOffsetValue = "-1"
// Constants for validating Successful reponse strings from server
const successfulPingResponse = "+PONG" + delimiter
const successfulResponse = "+OK" + delimiter

type RedisHost interface {
	Init()
	GetHostConfig() *model.HostConfig 
}