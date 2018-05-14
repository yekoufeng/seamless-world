package events

// IEvents 事件分发器接口
type IEvents interface {
	// 添加一个事件，并绑定到 objectInst 名为 callback的函数上，并返回句柄
	AddListener(evtName string, objInst interface{}, callback string) uint64
	// 把objinst上面绑的所有事件都删掉
	RemoveListenerByObjInst(objInst interface{})
	// 把某个句柄的事件去掉
	RemoveListenerByEvtHandle(evtHandle uint64)
	// 触发事件
	FireEvent(evtName string, args ...interface{})
	// 处理事件
	HandleEvent()
}
