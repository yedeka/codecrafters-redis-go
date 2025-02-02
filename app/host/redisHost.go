package host

import("github.com/codecrafters-io/redis-starter-go/app/model")

const linePrefix = "*"
const delimiter = "\r\n"
const lengthPrefix = "$"
const serverDelimiter = ":"
const pingCommand = "PING"

type RedisHost interface {
	Init()
	GetHostConfig() *model.HostConfig 
}