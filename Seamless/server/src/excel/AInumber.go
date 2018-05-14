package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AInumberData struct {
Id uint64
Score int
AInumber int
}

var AInumber map[uint64]AInumberData
var AInumberLock sync.RWMutex

func LoadAInumber(){
AInumberLock.Lock()
defer AInumberLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/AInumber.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &AInumber)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAInumberMap() map[uint64]AInumberData {
AInumberLock.RLock()
defer AInumberLock.RUnlock()

AInumber2 := make(map[uint64]AInumberData)
for k, v := range AInumber{
AInumber2[k] = v
}

return AInumber2
}

func GetAInumber(key uint64) (AInumberData, bool) {
AInumberLock.RLock()
defer AInumberLock.RUnlock()

val, ok := AInumber[key]

return val, ok
}

func GetAInumberMapLen() int {
return len(AInumber)
}

