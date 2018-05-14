package db

import (
	"strconv"
	"time"
	"zeus/dbservice"

	//"common"
	"encoding/json"
	"protoMsg"

	log "github.com/cihub/seelog"
)

// AnnuonceData 公告数据
type AnnuonceData struct {
	ID           uint64
	ServerID     uint32 // 服务器：微信（1），手Q（2）
	PlatID       uint8  // 平台：IOS（0），安卓（1）
	StartTime    uint32 // 开始时间
	EndTime      uint32 // 结束时间
	InternalTime uint32 // 滚动间隔时间 ：**秒/次
	Content      string // 公告内容
	Source       uint32 // 渠道号，由前端生成，不需要填写
	Serial       string // 流水号，由前端生成，不需要填写
}

type SliAnData []*AnnuonceData

func (a SliAnData) Len() int           { return len(a) }
func (a SliAnData) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SliAnData) Less(i, j int) bool { return a[i].ID < a[j].ID }

// AddAnnuonceData 添加公告数据
func AddAnnuonceData(data *AnnuonceData) bool {

	if data == nil {
		log.Warn("AddAnnuonceData is nil")
		return false
	} else {
		//log.Warn("AddAnnuonceData", data)
	}

	if hExists("AnnuonceData", data.ID) {
		return false
	}

	d, err := json.Marshal(*data)
	if err != nil {
		log.Info("AddAnnuonceData err ", err)
		return false
	}
	hSet("AnnuonceData", data.ID, string(d)) // 添加公告数据

	hSet("Annuoncing", data.ID, data.EndTime) // 进行中公告记录

	log.Info("AddAnnuonceData ", data)
	return true
}

// GetAnnuonceData 获取公告数据
func GetAnnuonceData(id uint64) *AnnuonceData {
	if !hExists("AnnuonceData", id) {
		return nil
	}

	v := hGet("AnnuonceData", id)
	var d *AnnuonceData
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GetAnnuonceData Failed to Unmarshal ", err)
		return nil
	}
	return d
}

// DelAnnuoncing 删除进行中公告(只进行中公告记录，不删数据)
func DelAnnuoncing(id uint64) bool {
	if !hExists("Annuoncing", id) {
		return false
	}

	hDEL("Annuoncing", id)

	data := GetAnnuonceData(id)
	if data == nil {
		return false
	}

	curtime := time.Now().Unix()
	if int64(data.EndTime) < curtime {
		return false
	}

	return true
}

// GetAllAnnuoncingData 获取所有进行中的公告数据
func GetAllAnnuoncingData() *protoMsg.InitAnnuonceInfoRet {

	retMsg := &protoMsg.InitAnnuonceInfoRet{}

	curtime := time.Now().Unix()
	idSet := hGetAll("Annuoncing")

	id := int64(0)
	endTime := int64(0)
	var err error = nil
	for idStr, endTimeStr := range idSet {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}

		endTime, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			continue
		}

		var data *AnnuonceData = GetAnnuonceData(uint64(id))
		if data == nil {
			continue
		}

		if endTime < curtime {
			DelAnnuoncing(uint64(id))
			continue
		}

		item := &protoMsg.AnnuonceInfo{
			Id:           data.ID,
			StartTime:    int64(data.StartTime),
			EndTime:      int64(data.EndTime),
			InternalTime: int64(data.InternalTime),
			Content:      data.Content,
		}

		retMsg.Item = append(retMsg.Item, item)
	}

	if len(retMsg.Item) == 0 {
		return nil
	}

	return retMsg
}

// PrintAnnuoncingData 打印所有进行中的公告数据
func PrintAnnuoncingData() {

	idSet := hGetAll("Annuoncing")
	log.Info("Annuoncing set ", idSet)
	for idStr, _ := range idSet {
		log.Info("Annuoncing id = ", idStr)

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Warn("Annuoncing data = nil")
			continue
		}

		var data *AnnuonceData = GetAnnuonceData(uint64(id))
		if data == nil {
			log.Warn("Annuoncing data = nil")
			continue
		}

		log.Info("Annuoncing data = ", data)
	}
}

// 获取全局公告id
func GetAnnuonceGlobalID() uint64 {
	result, err := dbservice.UIDGenerator().Get("annuceid")
	if err != nil {
		log.Warn("GetAnnuonceGlobalID failed ", err)
		return 0
	}

	return result
}

//  查询公告数据
func QueryAnnuonceData(beginTime uint32, endTime uint32) []*AnnuonceData {
	ret := make([]*AnnuonceData, 0)

	idMap := hGetAll("AnnuonceData")
	for _, s := range idMap {

		var d *AnnuonceData
		if err := json.Unmarshal([]byte(s), &d); err != nil {
			log.Warn("QueryAnnuonceData Failed to Unmarshal ", err)
			continue
		}

		if d.StartTime < beginTime || d.StartTime > endTime {
			continue
		}

		ret = append(ret, d)
	}

	return ret

}
