package db

import (
	"encoding/json"
	"fmt"
	"zeus/dbservice"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

const (
	teamPrefix      = "SpaceTeam"
	voiceInfoPrefix = "TeamVoice"
)

// TeamInfo 队伍信息
type TeamInfo struct {
	Id             uint64
	FixRatingIndex float32 // 获取FixRating的序列号.(使用rating/rating.go 中的GetFixRatingByIndex)
	MemList        []uint64
}

// SpaceTeamInfo 场景队伍信息
type SpaceTeamInfo struct {
	Id       uint64
	Teamtype uint32
	Teams    []TeamInfo
}

type spaceTeamUtil struct {
	uid uint64
}

type MemVoiceInfo struct {
	Uid      uint64 `json:"uid"`
	MemberId int32  `json:"memberid"`
}
type RedisVoiceInfo struct {
	MemList []*MemVoiceInfo `json:"memlist"`
}

// SpaceTeamUtil 队伍信息工具类
func SpaceTeamUtil(uid uint64) *spaceTeamUtil {
	return &spaceTeamUtil{
		uid: uid,
	}
}

func (r *spaceTeamUtil) key() string {
	return fmt.Sprintf("%s:%d", teamPrefix, r.uid)
}

// 保存队伍信息
func (r *spaceTeamUtil) SaveSpaceTeamInfo(info *SpaceTeamInfo) {

	if info == nil {
		log.Warn("SaveSpaceTeamInfo error info = ", info)
		return
	}

	d, e := json.Marshal(info)
	if e != nil {
		log.Warn("SaveSpaceTeamInfo error e = ", e)
	}
	if err := dbservice.CacheHSET(r.key(), r.uid, string(d)); err != nil {
		log.Error(err)
	}

	r.printfSpaceTeamInfo(info)
}

// 打印队伍信息
func (r *spaceTeamUtil) printfSpaceTeamInfo(info *SpaceTeamInfo) {
	if info == nil {
		return
	}

	log.Infof("打印队伍信息 sapce.dbid(%d), id(%d)", r.uid, info.Id)
	for _, teams := range info.Teams {
		log.Infof("队伍 id(%d)", teams.Id)
		log.Info("队伍成员", teams.MemList)
	}

}

// 获取队伍信息
func (r *spaceTeamUtil) GetSpaceTeamInfo() *SpaceTeamInfo {
	exists, err := dbservice.CacheHExists(r.key(), r.uid)
	if err != nil || !exists {
		log.Error("获取队伍信息失败 ", err)
		return nil
	}

	v, err := redis.String(dbservice.CacheHGET(r.key(), r.uid))
	if err != nil {
		log.Error("获取队伍信息失败 ", err)
		return nil
	}

	var d *SpaceTeamInfo
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GetSpaceTeamInfo Failed to Unmarshal ", err)
		return nil
	}
	return d
}

// 删除队伍信息
func (r *spaceTeamUtil) DelSpaceTeamInfo() {
	if err := dbservice.CacheHDEL(r.key(), r.uid); err != nil {
		log.Error("spaceid :", r.uid, " err:", err)
	}
}

// SetMatchSrvID 设置队伍所在的匹配服务器
func (r *spaceTeamUtil) SetMatchSrvID(srvID uint64) error {
	return dbservice.CacheHSET(r.key(), "matchsrvid", srvID)
}

// GetMatchSrvID 获取队伍所在的匹配服务器
func (r *spaceTeamUtil) GetMatchSrvID() (uint64, error) {
	return redis.Uint64(dbservice.CacheHGET(r.key(), "matchsrvid"))
}

// 获取新队伍全局唯一id
func GetTeamGlobalID() uint64 {
	result, err := dbservice.UIDGenerator().Get("teamglobalid")
	if err != nil {
		log.Error("GetTeamGlobalID failed ", err)
		return 0
	}

	return result
}

func GetTeamVoiceInfo(teamId uint64) []*MemVoiceInfo {
	c := dbservice.GetServerRedis()
	defer c.Close()

	key := fmt.Sprintf("%s:%d", voiceInfoPrefix, teamId)
	data, err := redis.Bytes(c.Do("get", key))
	if err != nil {
		log.Error(err)
		return nil
	}

	redisInfo := &RedisVoiceInfo{}
	err = json.Unmarshal(data, redisInfo)
	if err != nil {
		log.Error(err)
		return nil
	}

	return redisInfo.MemList
}

// SetTeamVoiceInfo TODO:什么时候删除是个问题，需要解决
func SetTeamVoiceInfo(teamId uint64, memlist []*MemVoiceInfo) {
	redisInfo := &RedisVoiceInfo{
		MemList: memlist,
	}
	data, err := json.Marshal(redisInfo)
	if err != nil {
		log.Error(err)
		return
	}
	c := dbservice.GetServerRedis()
	defer c.Close()

	key := fmt.Sprintf("%s:%d", voiceInfoPrefix, teamId)

	c.Do("set", key, data)
}

func DelTeamVoiceInfo(teamId uint64) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	key := fmt.Sprintf("%s:%d", voiceInfoPrefix, teamId)

	c.Do("del", key)
}
