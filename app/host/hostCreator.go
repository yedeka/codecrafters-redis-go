package host

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"strconv"
	
	"github.com/codecrafters-io/redis-starter-go/app/model"
)

func CreateHost(listeningPort string, 
	replicationData string, 
	hostAddress string) {
		if replicationData != "leader" {
			fmt.Println("Creating follower")
			follower, err := createFollower(replicationData)
			if nil != err {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			follower.Init()
			createReplicaMaster(hostAddress, listeningPort, follower.GetHostConfig()).Init()
		} else {
			createMaster(hostAddress, listeningPort).Init()
		}
		 
}

func createReplicaMaster(hostAddress string, listeningPort string, followerHostConfig *model.HostConfig) RedisHost {
	return Master {
		hostIpAddress: hostAddress,
		listeningPort: listeningPort,
		hostConfig: followerHostConfig,
	}	
}

func createMaster(hostAddress string, listeningPort string) RedisHost {
		return Master {
			hostIpAddress: hostAddress,
			listeningPort: listeningPort,
			hostConfig: &model.HostConfig{
				IsMaster: true,
			},
		}	
}

func createFollower(replicationData string) (RedisHost, error) {
	replicationFollowerData := strings.Split(replicationData, " ")
			masterPort, err := strconv.Atoi(replicationFollowerData[1])
			if nil != err {
				return nil, errors.New("non numeric port configuration passed for server")
			}
			return Follower {
				hostConfig: &model.HostConfig{
					IsMaster: false,
					MasterProps: model.MasterConfig{
						Host: replicationFollowerData[0],
						Port: masterPort,
					},
				},
			}, nil
}