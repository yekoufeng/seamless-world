package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type StoreData struct {
Id uint64
Type uint64
Name string
Price uint64
RelationID uint64
State uint64
}

var store map[uint64]StoreData
var storeLock sync.RWMutex

func LoadStore(){
storeLock.Lock()
defer storeLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/store.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &store)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetStoreMap() map[uint64]StoreData {
storeLock.RLock()
defer storeLock.RUnlock()

store2 := make(map[uint64]StoreData)
for k, v := range store{
store2[k] = v
}

return store2
}

func GetStore(key uint64) (StoreData, bool) {
storeLock.RLock()
defer storeLock.RUnlock()

val, ok := store[key]

return val, ok
}

func GetStoreMapLen() int {
return len(store)
}

