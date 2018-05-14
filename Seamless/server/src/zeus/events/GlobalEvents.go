package events

import (
	"fmt"
	"reflect"
	"sync"
	"zeus/dbservice"
	"zeus/serializer"

	log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

// GlobalEvents 全局事件分发器
// 使用Redis的事件分布机制，实现分布式系统的全局事件定阅机制
type GlobalEvents struct {
	evtNameMap map[string][]uint64
	evtObjMap  map[reflect.Value][]uint64
	evtMap     map[uint64]reflect.Value

	mtx *sync.Mutex

	psc    redis.PubSubConn
	inited bool
	closed bool
	eventQ chan *redis.Message
}

// NewGlobalEventsInst 新建全局Events对象
func NewGlobalEventsInst() *GlobalEvents {
	globalEventsInst := &GlobalEvents{}
	globalEventsInst.evtNameMap = make(map[string][]uint64)
	globalEventsInst.evtObjMap = make(map[reflect.Value][]uint64)
	globalEventsInst.evtMap = make(map[uint64]reflect.Value)

	globalEventsInst.mtx = &sync.Mutex{}

	c, err := dbservice.GetSingletonConn()
	if err != nil {
		panic(err)
	}
	globalEventsInst.psc = redis.PubSubConn{Conn: c}
	globalEventsInst.eventQ = make(chan *redis.Message, 2000)
	globalEventsInst.inited = false
	globalEventsInst.closed = false

	return globalEventsInst
}

// Destroy 销毁事件管理器
func (evt *GlobalEvents) Destroy() {
	if evt.closed {
		return
	}
	evt.closed = true
	evt.evtNameMap = nil
	evt.evtObjMap = nil
	evt.evtMap = nil
	evt.psc.Unsubscribe()
	evt.psc.Close()
}

// AddListener 增加一个事件监听
func (evt *GlobalEvents) AddListener(evtName string, objInst interface{}, callback string) uint64 {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err, evtName, callback)
			if viper.GetString("Config.Recover") == "0" {
				panic(fmt.Sprintln(err, evtName, callback))
			}
		}
	}()

	evt.mtx.Lock()
	defer evt.mtx.Unlock()

	if _, ok := evt.evtNameMap[evtName]; !ok {
		if err := evt.psc.Subscribe(evtName); err != nil {
			return 0
		}

		if evt.inited == false {
			go func() {
				for {
					switch n := evt.psc.Receive().(type) {
					case redis.Message:
						evt.eventQ <- &n
					case error:
						log.Error(n)
						evt.Destroy()
						return
					}
				}
			}()
			evt.inited = true
		}
	}

	v := reflect.ValueOf(objInst)
	m := v.MethodByName(callback)
	for k, exitsM := range evt.evtMap {
		if m == exitsM {
			return k
		}
	}

	id, err := dbservice.UIDGenerator().Get("handler")
	if err != nil {
		return 0
	}
	evt.evtMap[id] = m
	evt.evtNameMap[evtName] = append(evt.evtNameMap[evtName], id)
	evt.evtObjMap[v] = append(evt.evtObjMap[v], id)
	return id
}

// RemoveListenerByObjInst 删除该实例下所有事件处理方法
func (evt *GlobalEvents) RemoveListenerByObjInst(objInst interface{}) {

	evt.mtx.Lock()
	defer evt.mtx.Unlock()

	v := reflect.ValueOf(objInst)
	if l, ok := evt.evtObjMap[v]; ok {
		for _, h := range l {
			evt.doRemoveFromEvtMap(h)
		}
		delete(evt.evtObjMap, v)
	}
}

// RemoveListenerByEvtHandle 删除事件处理方法
func (evt *GlobalEvents) RemoveListenerByEvtHandle(evtHandle uint64) {
	evt.mtx.Lock()
	defer evt.mtx.Unlock()

	evt.doRemoveFromEvtMap(evtHandle)
	evt.doRemoveFromObjMap(evtHandle)
}

// FireEvent 触发一个事件
func (evt *GlobalEvents) FireEvent(evtName string, args ...interface{}) {
	c := dbservice.GetSingletonRedis()
	defer c.Close()

	data := serializer.Serialize(args...)
	_, err := c.Do("Publish", evtName, data)
	if err != nil {
		log.Error(err)
	}
}

// HandleEvent 处理一个事件
func (evt *GlobalEvents) HandleEvent() {

	evt.mtx.Lock()
	defer evt.mtx.Unlock()

	handled := 0
	var msg *redis.Message

	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			if viper.GetString("Config.Recover") == "0" {
				panic(fmt.Sprintln(msg))
			}
		}
	}()

	for {
		select {
		case msg = <-evt.eventQ:
			if _, ok := evt.evtNameMap[msg.Channel]; ok {
				if len(evt.evtNameMap[msg.Channel]) == 0 {
					delete(evt.evtNameMap, msg.Channel)
					evt.psc.Unsubscribe(msg.Channel)
				}

				for k := 0; k < len(evt.evtNameMap[msg.Channel]); {
					h := evt.evtNameMap[msg.Channel][k]
					if f, ok := evt.evtMap[h]; ok {
						args := serializer.UnSerialize(msg.Data)
						callArgs := []reflect.Value{}
						for _, arg := range args {
							callArgs = append(callArgs, reflect.ValueOf(arg))
						}
						f.Call(callArgs)
						k++
					} else {
						evt.evtNameMap[msg.Channel] = append(evt.evtNameMap[msg.Channel][:k], evt.evtNameMap[msg.Channel][k+1:]...)
					}
				}
			}
			handled++
		default:
			return
		}

		if handled >= 10 {
			return
		}
	}
}

func (evt *GlobalEvents) doRemoveFromEvtMap(evtHandle uint64) {
	if _, ok := evt.evtMap[evtHandle]; ok {
		delete(evt.evtMap, evtHandle)
	}
}

func (evt *GlobalEvents) doRemoveFromObjMap(evtHandle uint64) {
	for k, lt := range evt.evtObjMap {
		for i, v := range lt {
			if v == evtHandle {
				evt.evtObjMap[k] = append(evt.evtObjMap[k][:i], evt.evtObjMap[k][i+1:]...)
				return
			}
		}
	}
}
