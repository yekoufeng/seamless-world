package server

import (
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"
	"zeus/dbservice"

	log "github.com/cihub/seelog"
)

// GetValidSrvPort 根据传入的最小端口和最大端口号，找一个可用端口，返回一个可用来绑定的地址
func GetValidSrvPort(minPort, maxPort int) string {
	if maxPort <= minPort {
		return strconv.Itoa(minPort)
	}
	maxRetry := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		port := r.Intn(int(maxPort-minPort)) + int(minPort)
		portstr := strconv.Itoa(port)

		l, err := net.Listen("tcp", ":"+portstr)
		if err != nil {
			log.Error("尝试监听端口失败 ", err)
			if maxRetry > 10 {
				return ""
			}
			maxRetry++
			time.Sleep(3 * time.Second)
			continue
		} else {
			err = l.Close()
			if err != nil {
				log.Error("关闭监听端口失败", err)
			}
			return portstr
		}
	}
}

/*
// GetEntityTempID 获取实体唯一的临时ID
func GetEntityTempID() uint64 {
	u, err := dbservice.UIDGenerator().Get("entity")
	if err != nil {
		return math.MaxUint64
	}
	return u
}
*/

// CreateNewUID 生成新的UID
func CreateNewUID() (uint64, error) {
	u, err := dbservice.UIDGenerator().Get("user")
	if err != nil {
		return math.MaxUint64, err
	}
	return u, nil
}
