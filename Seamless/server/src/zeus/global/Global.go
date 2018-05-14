package global

import (
	"encoding/json"
	"errors"
	"zeus/dbservice"
	"zeus/iserver"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

const (
	globalVariantTypeString      = 1
	globalVariantTypeInt         = 2
	globalVariantTypeEntityProxy = 3
	globalVariantTypeIntSlice    = 4
)

type globalInfo struct {
	variantType uint8
	val         interface{}
}

func newGlobalInfo(varType uint8, val interface{}) *globalInfo {
	return &globalInfo{
		variantType: varType,
		val:         val,
	}
}

// Global 支持分布式结构的全局数据结构
type Global struct {
	variants map[string]*globalInfo
}

var globalInst *Global

// GetGlobalInst 获取到全局变量指针
func GetGlobalInst() *Global {

	if globalInst == nil {
		globalInst = &Global{}
		globalInst.init()
	}

	return globalInst
}

func (g *Global) init() {
	iserver.GetSrvInst().AddListener(g.getGlobalRefreshEvtName(), g, "GlobalVariantChanged")
	g.refreshGlobal()
}

func (g *Global) getGlobalRefreshEvtName() string {
	return "Event_GlobalRefresh"
}

func (g *Global) getGlobalKeyPrefix() string {
	return "Global_Variants:"
}

func (g *Global) getGlobalKeysKey() string {
	return "Global_Keys"
}

func (g *Global) refreshGlobal() {
	//从Redis的Global表中获取全局变量
	c := dbservice.GetServerRedis()
	defer c.Close()

	values, err := redis.Values(c.Do("SMEMBERS", g.getGlobalKeysKey()))
	if err != nil {
		log.Error("获取全局变量异常")
		return
	}

	g.variants = make(map[string]*globalInfo)

	for i := 0; i < len(values); i++ {
		globalName, _ := redis.String(values[i], nil)
		variantType, val, err := g.getGlobalFromDB(globalName)

		if err != nil {
			log.Error("全局变量刷新异常")
			continue
		}

		g.variants[globalName] = newGlobalInfo(variantType, val)
	}

}

func (g *Global) getGlobalFromDB(globalName string) (variantType uint8, val interface{}, err error) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	keyName := g.getGlobalKeyPrefix() + globalName
	values, e := redis.Values(c.Do("HGETALL", keyName))
	if e != nil {
		err = e
		return
	}

	vars := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		str, _ := redis.String(values[i], nil)
		vars[str] = values[i+1]
	}

	_, ok := vars["variantType"]
	if !ok {
		err = errors.New("全局变量表异常")
		return
	}

	err = nil
	varType, _ := redis.Int(vars["variantType"], nil)
	variantType = uint8(varType)

	switch variantType {
	case globalVariantTypeString:
		val, _ = redis.String(vars["stringVal"], nil)
	case globalVariantTypeInt:
		val, _ = redis.Int(vars["intVal"], nil)
	case globalVariantTypeEntityProxy:
		e := make(map[string]uint64)
		e["entityID"], _ = redis.Uint64(vars["entityID"], nil)
		e["srvID"], _ = redis.Uint64(vars["srvID"], nil)
		e["cellID"], _ = redis.Uint64(vars["cellID"], nil)
		val = e
	case globalVariantTypeIntSlice:
		var data []byte
		data, err = redis.Bytes(vars["intSliceVal"], nil)
		if err != nil {
			return
		}

		var ret []int
		err = json.Unmarshal(data, &ret)
		if err != nil {
			return
		}
		val = ret
	default:
		err = errors.New("错误的全局变量类型")
		return
	}

	return
}

// GlobalVariantChanged 当全局变量发生改变时收到通知
func (g *Global) GlobalVariantChanged(args string) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	keyName := g.getGlobalKeyPrefix() + string(args)

	isKeyExisted, err := redis.Bool(c.Do("EXISTS", keyName))
	if err != nil {
		log.Error(err)
		return
	}

	if isKeyExisted {

		variantType, val, err := g.getGlobalFromDB(args)

		if err != nil {
			log.Error(err)
			return
		}

		g.variants[args] = newGlobalInfo(variantType, val)
	} else {
		delete(g.variants, args)
	}

}

// RemoveGlobal 移除全局变量
func (g *Global) RemoveGlobal(globalName string) {

	if _, ok := g.variants[globalName]; !ok {
		return
	}

	c := dbservice.GetServerRedis()
	defer c.Close()

	c.Do("DEL", g.getGlobalKeyPrefix()+globalName)
	c.Do("SREM", g.getGlobalKeysKey(), globalName)

	iserver.GetSrvInst().FireEvent(g.getGlobalRefreshEvtName(), globalName)
}

// SetGlobalStr 设置全局变量字符串
func (g *Global) SetGlobalStr(globalName string, val string) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	var hashArgs []interface{}
	hashArgs = append(hashArgs, g.getGlobalKeyPrefix()+globalName)
	hashArgs = append(hashArgs, "variantType")
	hashArgs = append(hashArgs, globalVariantTypeString)
	hashArgs = append(hashArgs, "stringVal")
	hashArgs = append(hashArgs, val)
	c.Do("HMSET", hashArgs...)

	// c.Do("HSET", g.getGlobalKeyPrefix()+globalName, "variantType", globalVariantTypeString)
	// c.Do("HSET", g.getGlobalKeyPrefix()+globalName, "stringVal", val)

	c.Do("SADD", g.getGlobalKeysKey(), globalName)

	iserver.GetSrvInst().FireEvent(g.getGlobalRefreshEvtName(), globalName)
}

// SetGlobalInt 设置全局变量Int
func (g *Global) SetGlobalInt(globalName string, val int) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	var hashArgs []interface{}
	hashArgs = append(hashArgs, g.getGlobalKeyPrefix()+globalName)
	hashArgs = append(hashArgs, "variantType")
	hashArgs = append(hashArgs, globalVariantTypeInt)
	hashArgs = append(hashArgs, "intVal")
	hashArgs = append(hashArgs, val)
	c.Do("HMSET", hashArgs...)

	// c.Do("HSET", g.getGlobalKeyPrefix()+globalName, "variantType", globalVariantTypeInt)
	// c.Do("HSET", g.getGlobalKeyPrefix()+globalName, "intVal", val)

	c.Do("SADD", g.getGlobalKeysKey(), globalName)

	iserver.GetSrvInst().FireEvent(g.getGlobalRefreshEvtName(), globalName)
}

// SetGlobalEntityProxy 设置全局实体代理对象
func (g *Global) SetGlobalEntityProxy(globalName string, entityID, srvID, cellID uint64) {
	c := dbservice.GetServerRedis()
	defer c.Close()

	var hashArgs []interface{}
	hashArgs = append(hashArgs, g.getGlobalKeyPrefix()+globalName)
	hashArgs = append(hashArgs, "variantType")
	hashArgs = append(hashArgs, globalVariantTypeEntityProxy)
	hashArgs = append(hashArgs, "entityID")
	hashArgs = append(hashArgs, entityID)
	hashArgs = append(hashArgs, "srvID")
	hashArgs = append(hashArgs, srvID)
	hashArgs = append(hashArgs, "cellID")
	hashArgs = append(hashArgs, cellID)

	c.Do("HMSET", hashArgs...)
	c.Do("SADD", g.getGlobalKeysKey(), globalName)

	iserver.GetSrvInst().FireEvent(g.getGlobalRefreshEvtName(), globalName)
}

// SetGlobalIntSlice 设置int slice
func (g *Global) SetGlobalIntSlice(globalName string, val []int) {
	data, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return
	}

	c := dbservice.GetServerRedis()
	defer c.Close()

	var hashArgs []interface{}
	hashArgs = append(hashArgs, g.getGlobalKeyPrefix()+globalName)
	hashArgs = append(hashArgs, "variantType")
	hashArgs = append(hashArgs, globalVariantTypeIntSlice)
	hashArgs = append(hashArgs, "intSliceVal")
	hashArgs = append(hashArgs, data)

	c.Do("HMSET", hashArgs...)
	c.Do("SADD", g.getGlobalKeysKey(), globalName)

	iserver.GetSrvInst().FireEvent(g.getGlobalRefreshEvtName(), globalName)
}

// GetGlobalStr 获取全局变量字符串
func (g *Global) GetGlobalStr(globalName string) string {

	info, ok := g.variants[globalName]
	if !ok {
		return ""
	}

	if info.variantType != globalVariantTypeString {
		log.Warn("全局变量类型错误")
		return ""
	}

	return info.val.(string)
}

// GetGlobalInt 获取全局变量Int
func (g *Global) GetGlobalInt(globalName string) int {

	info, ok := g.variants[globalName]
	if !ok {
		return 0
	}

	if info.variantType != globalVariantTypeInt {
		log.Warn("全局变量类型错误")
		return 0
	}

	return info.val.(int)
}

// GetGlobalEntityProxy 获取全局变量Entity代理
func (g *Global) GetGlobalEntityProxy(globalName string) (uint64, uint64, uint64) {
	info, ok := g.variants[globalName]
	if !ok {
		return 0, 0, 0
	}

	if info.variantType != globalVariantTypeEntityProxy {
		log.Warn("全局变量类型错误")
		return 0, 0, 0
	}

	e := info.val.(map[string]uint64)
	return e["entityID"], e["srvID"], e["cellID"]
}

// GetGlobalIntSlice 获取int slice
func (g *Global) GetGlobalIntSlice(globalName string) []int {
	info, ok := g.variants[globalName]
	if !ok {
		return nil
	}

	if info.variantType != globalVariantTypeIntSlice {
		log.Warn("全局变量类型错误 ", globalName, " 类型 ", info.variantType)
		return nil
	}

	return info.val.([]int)
}
