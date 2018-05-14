package common

import (
	"excel"
	"math"
	"strconv"
	"time"
	"zeus/dbservice"
	"zeus/linmath"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

func StringToUint64(i string) uint64 {
	d, e := strconv.ParseUint(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return d
}
func StringToInt64(i string) int64 {
	d, e := strconv.ParseInt(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return d
}
func StringToUint32(i string) uint32 {
	d, e := strconv.ParseUint(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return uint32(d)
}
func StringToInt(i string) int {
	d, e := strconv.ParseInt(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return int(d)
}
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func InitMsg() {
	def := msgdef.GetMsgDef()
	for k, v := range ProtoMap {
		def.RegMsg(k, v)
	}
}

func Distance(v linmath.Vector3, o linmath.Vector3) float32 {
	dx := v.X - o.X
	dz := v.Z - o.Z

	return float32(math.Sqrt(float64(dx*dx + dz*dz)))
}

func CreateNewMailID() uint64 {
	u, err := dbservice.UIDGenerator().Get("mail")
	if err != nil {
		return 0
	}
	return u
}

func GetDefiniteequinox(x1 float32, x2 float32, rate float32) float32 {
	return (x1 + rate*x2) / (1 + rate)
}

func GetTBSystemValue(id uint64) uint {
	base, ok := excel.GetSystem(id)
	if !ok {
		return 0
	}

	return uint(base.Value)
}

func GetUsers() []uint64 {
	c := dbservice.Get()
	defer c.Close()

	reply, err := c.Do("HVALS", "accounts")
	if err != nil {
		return nil
	}
	if reply == nil {
		return nil
	}

	values, err := redis.Values(reply, nil)
	if err != nil {
		return nil
	}

	result := make([]uint64, 0, len(values))
	for _, v := range values {
		uid, err := redis.Uint64(v, nil)
		if err != nil {
			continue
		}

		result = append(result, uid)
	}

	log.Info("获取用户总数", len(result))
	return result
}

func GetSeason() int {
	season := excel.GetSeasonMap()

	var tmpID = 0
	var tmpTime int64 = 0
	for _, v := range season {
		tm, _ := time.Parse("2006|01|02", v.EndTime)
		if time.Now().Unix() <= tm.Unix() {
			if tm.Unix() < tmpTime || tmpTime == 0 {
				tmpTime = tm.Unix()
				tmpID = int(v.Id)
			}
		}
	}

	if tmpTime == 0 {
		log.Warn("赛季配置表出问题")
	}
	return tmpID
}
