package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type BarrierData struct {
Id uint64
Type uint64
Name string
HitSpeed float32
}

var barrier map[uint64]BarrierData
var barrierLock sync.RWMutex

func LoadBarrier(){
barrierLock.Lock()
defer barrierLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/barrier.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &barrier)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetBarrierMap() map[uint64]BarrierData {
barrierLock.RLock()
defer barrierLock.RUnlock()

barrier2 := make(map[uint64]BarrierData)
for k, v := range barrier{
barrier2[k] = v
}

return barrier2
}

func GetBarrier(key uint64) (BarrierData, bool) {
barrierLock.RLock()
defer barrierLock.RUnlock()

val, ok := barrier[key]

return val, ok
}

func GetBarrierMapLen() int {
return len(barrier)
}

