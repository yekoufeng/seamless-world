package db

import (
	"encoding/json"
	"fmt"
	"zeus/linmath"

	log "github.com/cihub/seelog"
)

// 玩家的地图副本相关数据

const (
	playerMapPrefix = "PlayerMap"
)

const (
	fieldMapPos = "MapPos"
)

type playerMapUtil struct {
	uid uint64
}

// PlayerMapUtil 获取工具类
func PlayerMapUtil(uid uint64) *playerMapUtil {
	return &playerMapUtil{
		uid: uid,
	}
}

// PlayerMapData 地图相关数据
type PlayerMapData struct {
	ID      uint64
	MapID   uint64          // 地图ID
	MapName string          // 地图名
	Pos     linmath.Vector3 //坐标
	Rota    linmath.Vector3 //旋转
}

// SetMapData 添加公告数据
func (u *playerMapUtil) SetPlayerMapData(data *PlayerMapData) bool {
	if data == nil {
		log.Warn("SetPlayerMapData is nil")
		return false
	} else {
		//log.Warn("AddAnnuonceData", data)
	}

	d, err := json.Marshal(*data)
	if err != nil {
		log.Info("SetPlayerMapData err ", err)
		return false
	}
	hSet(u.key(), fieldMapPos, string(d)) // 添加公告数据

	return true
}

// GetPlayerMapData 获取地图数据
func (u *playerMapUtil) GetPlayerMapData() *PlayerMapData {
	if !hExists(u.key(), fieldMapPos) {
		return nil
	}

	v := hGet(u.key(), fieldMapPos)
	var d *PlayerMapData
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GePlayerMapData Failed to Unmarshal ", err)
		return nil
	}
	return d
}

func (u *playerMapUtil) key() string {
	return fmt.Sprintf("%s:%d", playerTempPrefix, u.uid)
}
