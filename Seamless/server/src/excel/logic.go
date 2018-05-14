package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type LogicData struct {
Id uint64
Prestate int
Nextstate int
Stiffness float32
}

var logic map[uint64]LogicData
var logicLock sync.RWMutex

func LoadLogic(){
logicLock.Lock()
defer logicLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/logic.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &logic)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetLogicMap() map[uint64]LogicData {
logicLock.RLock()
defer logicLock.RUnlock()

logic2 := make(map[uint64]LogicData)
for k, v := range logic{
logic2[k] = v
}

return logic2
}

func GetLogic(key uint64) (LogicData, bool) {
logicLock.RLock()
defer logicLock.RUnlock()

val, ok := logic[key]

return val, ok
}

func GetLogicMapLen() int {
return len(logic)
}

