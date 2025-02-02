package host

import("github.com/codecrafters-io/redis-starter-go/app/model")

const linePrefix = "*"
const delimiter = "\r\n"
const lengthPrefix = "$"
const serverDelimiter = ":"
const pingCommand = "PING"
const replConfCommand = "REPLCONF"

const listeningPortConfKey = "listening-port"
const replicationPort = "6380"
const capacityKey = "capa"
const defaultCapacityValue = "psync2"

const successfulPingResponse = "+PONG" + delimiter
const successfulResponse = "+OK" + delimiter

type RedisHost interface {
	Init()
	GetHostConfig() *model.HostConfig 
}