package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AImatchData struct {
Id uint64
Name string
Value float32
}

var AImatch map[uint64]AImatchData
var AImatchLock sync.RWMutex

func LoadAImatch(){
AImatchLock.Lock()
defer AImatchLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/AImatch.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &AImatch)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAImatchMap() map[uint64]AImatchData {
AImatchLock.RLock()
defer AImatchLock.RUnlock()

AImatch2 := make(map[uint64]AImatchData)
for k, v := range AImatch{
AImatch2[k] = v
}

return AImatch2
}

func GetAImatch(key uint64) (AImatchData, bool) {
AImatchLock.RLock()
defer AImatchLock.RUnlock()

val, ok := AImatch[key]

return val, ok
}

func GetAImatchMapLen() int {
return len(AImatch)
}

