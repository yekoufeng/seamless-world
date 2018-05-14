package serverMgr

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"zeus/dbservice"
	"zeus/iserver"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	serverMapRW sync.RWMutex
	serverMap   map[uint8]iserver.ServerList
	succCount   uint32
	totalCount  uint32
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer() *LoadBalancer {
	app := &LoadBalancer{}
	if _, err := app.GetServerList(); err != nil {
		log.Error("Init failed", err)
		return nil
	}
	return app
}

// ErrServerBusy 服务器忙
var ErrServerBusy = errors.New("Get server failed, busy")

// GetServerByType 轮询获取服务器
func (app *LoadBalancer) GetServerByType(t uint8) (*iserver.ServerInfo, error) {
	// 当成功获取服务器的次数超过指定次数时, 重新刷新服务器列表
	if app.succCount >= 50 {
		if _, err := app.GetServerList(); err != nil {
			log.Error("Refresh server list failed", err)
		}
		atomic.StoreUint32(&app.succCount, 0)
	}

	app.serverMapRW.RLock()
	defer app.serverMapRW.RUnlock()

	if list, ok := app.serverMap[t]; ok {
		length := list.Len()
		if length <= 0 {
			return nil, fmt.Errorf("Server list is empty, Type %d", t)
		}

		atomic.AddUint32(&app.totalCount, 1)
		atomic.AddUint32(&app.succCount, 1)
		index := app.totalCount % uint32(length)

		tryTime := 0
		maxLoad := viper.GetInt("Config.MaxLoad")
		if maxLoad == 0 {
			maxLoad = 50
		}
		for {
			if tryTime > length {
				return nil, ErrServerBusy
			}
			tryTime++

			srv := list[index%uint32(length)]
			if srv.Load < maxLoad {
				return srv, nil
			}

			index++
		}
	}
	return nil, fmt.Errorf("Cant get server, Type %d not existed", t)
}

// GetServerList 获取最新的服务器列表
func (app *LoadBalancer) GetServerList() ([]*iserver.ServerInfo, error) {
	var list []*iserver.ServerInfo
	if err := dbservice.GetServerList(&list); err != nil {
		return nil, err
	}

	app.serverMapRW.Lock()
	defer app.serverMapRW.Unlock()

	app.serverMap = make(map[uint8]iserver.ServerList)
	for _, v := range list {
		app.serverMap[v.Type] = append(app.serverMap[v.Type], v)
	}

	atomic.StoreUint32(&app.succCount, 0)
	return list, nil
}
