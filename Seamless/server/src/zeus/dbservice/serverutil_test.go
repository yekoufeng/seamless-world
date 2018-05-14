package dbservice

import (
	"testing"
)

type ServerInfo struct {
	ServerID     uint64 `redis:"serverid"`
	Type         uint8  `redis:"type"`
	OuterAddress string `redis:"outeraddr"`
	InnerAddress string `redis:"inneraddr"`
	Console      uint64 `redis:"console"`
	Load         int    `redis:"load"`
	Token        string `redis:"token"`
	Status       int    `redis:"status"`
}

func TestServerUtilGetServerList(t *testing.T) {
	//c, _ := redis.Dial("tcp", "192.168.150.190:6379")
	//defer c.Close()

	// var list []*ServerInfo

	// err := GetServerList(c, &list)
	// fmt.Println(err)
	// fmt.Println(list[0])

	// util := new(serverUtil)
	// util.uid = 10000
	// util.c = c
	// fmt.Println(util.GetToken())

	//util = nil
}
