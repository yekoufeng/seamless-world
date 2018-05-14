package safecontainer

import "sync"

// NewSafeMap 新建一个安全map
func NewSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[interface{}]interface{})
	return sm

}

// SafeMap 线程安全的map，使用读写锁保护，写时阻塞，读时共享
type SafeMap struct {
	sync.RWMutex
	Map map[interface{}]interface{}
}

// Get 获取
func (sm *SafeMap) Get(key interface{}) interface{} {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

// Set 设置
func (sm *SafeMap) Set(key interface{}, value interface{}) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

// Remove 删除
func (sm *SafeMap) Remove(key interface{}) {
	sm.Lock()
	delete(sm.Map, key)
	sm.Unlock()
}

// IsExist 是否存在键值
func (sm *SafeMap) IsExist(key interface{}) bool {
	sm.RLock()
	_, ok := sm.Map[key]
	sm.RUnlock()

	return ok
}

// Travsal 遍历
func (sm *SafeMap) Travsal(cb func(interface{}, interface{})) {
	sm.RLock()

	for k, e := range sm.Map {
		cb(k, e)
	}

	sm.RUnlock()
}
