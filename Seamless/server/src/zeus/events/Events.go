package events

/*
	由于服务器是多线程环境，所以事件分发机制意义不大
*/

// // Events 本地事件分发器
// type Events struct {
// 	evtNameMap map[string][]uint64
// 	evtObjMap  map[reflect.Value][]uint64
// 	evtMap     map[uint64]reflect.Value

// 	eventID uint64
// }

// // GetEventsInst 获取本地Events对象
// // func GetEventsInst() IEvents {
// // 	if eventsInst == nil {
// // 		eventsInst = new(Events)
// // 		eventsInst.evtNameMap = make(map[string][]uint64)
// // 		eventsInst.evtObjMap = make(map[reflect.Value][]uint64)
// // 		eventsInst.evtMap = make(map[uint64]reflect.Value)

// // 		eventID = 1
// // 	}

// // 	return eventsInst
// // }

// // NewEventsInst 新建本地Events对象
// func NewEventsInst() *Events {
// 	eventsInst := &Events{}
// 	eventsInst.evtNameMap = make(map[string][]uint64)
// 	eventsInst.evtObjMap = make(map[reflect.Value][]uint64)
// 	eventsInst.evtMap = make(map[uint64]reflect.Value)
// 	eventsInst.eventID = 1

// 	return eventsInst
// }

// // AddListener 增加一个事件监听
// func (evt *Events) AddListener(evtName string, objInst interface{}, callback string) uint64 {
// 	id := evt.eventID

// 	v := reflect.ValueOf(objInst)
// 	m := v.MethodByName(callback)
// 	for k, exitsM := range evt.evtMap {
// 		if m == exitsM {
// 			return k
// 		}
// 	}

// 	evt.evtMap[id] = m
// 	evt.evtNameMap[evtName] = append(evt.evtNameMap[evtName], id)
// 	evt.evtObjMap[v] = append(evt.evtObjMap[v], id)

// 	evt.eventID++
// 	return id
// }

// // RemoveListenerByObjInst 删除该实例下所有事件处理方法
// func (evt *Events) RemoveListenerByObjInst(objInst interface{}) {
// 	v := reflect.ValueOf(objInst)
// 	if l, ok := evt.evtObjMap[v]; ok {
// 		for _, h := range l {
// 			evt.doRemoveFromEvtMap(h)
// 		}
// 		delete(evt.evtObjMap, v)
// 	}
// }

// // RemoveListenerByEvtHandle 删除事件处理方法
// func (evt *Events) RemoveListenerByEvtHandle(evtHandle uint64) {
// 	evt.doRemoveFromEvtMap(evtHandle)
// 	evt.doRemoveFromObjMap(evtHandle)
// }

// // FireEvent 触发一个事件
// func (evt *Events) FireEvent(evtName string, args ...interface{}) {
// 	if _, ok := evt.evtNameMap[evtName]; ok {
// 		if len(evt.evtNameMap[evtName]) == 0 {
// 			delete(evt.evtNameMap, evtName)
// 		}

// 		for k := 0; k < len(evt.evtNameMap[evtName]); {
// 			h := evt.evtNameMap[evtName][k]
// 			if f, ok := evt.evtMap[h]; ok {
// 				callArgs := []reflect.Value{}
// 				for _, arg := range args {
// 					callArgs = append(callArgs, reflect.ValueOf(arg))
// 				}
// 				f.Call(callArgs)
// 				k++
// 			} else {
// 				evt.evtNameMap[evtName] = append(evt.evtNameMap[evtName][:k], evt.evtNameMap[evtName][k+1:]...)
// 			}
// 		}
// 	}
// }

// // HandleEvent 处理事件
// func (evt *Events) HandleEvent() {

// }

// func (evt *Events) doRemoveFromEvtMap(evtHandle uint64) {
// 	if _, ok := evt.evtMap[evtHandle]; ok {
// 		delete(evt.evtMap, evtHandle)
// 	}
// }

// func (evt *Events) doRemoveFromObjMap(evtHandle uint64) {
// 	for k, lt := range evt.evtObjMap {
// 		for i, v := range lt {
// 			if v == evtHandle {
// 				evt.evtObjMap[k] = append(evt.evtObjMap[k][:i], evt.evtObjMap[k][i+1:]...)
// 				return
// 			}
// 		}
// 	}
// }
