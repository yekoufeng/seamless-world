package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type MeterfixData struct {
Id uint64
Dist uint64
AttackAgain float32
BaseAcc float32
BaseBeDetectedRate float32
}

var meterfix map[uint64]MeterfixData
var meterfixLock sync.RWMutex

func LoadMeterfix(){
meterfixLock.Lock()
defer meterfixLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/meterfix.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &meterfix)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetMeterfixMap() map[uint64]MeterfixData {
meterfixLock.RLock()
defer meterfixLock.RUnlock()

meterfix2 := make(map[uint64]MeterfixData)
for k, v := range meterfix{
meterfix2[k] = v
}

return meterfix2
}

func GetMeterfix(key uint64) (MeterfixData, bool) {
meterfixLock.RLock()
defer meterfixLock.RUnlock()

val, ok := meterfix[key]

return val, ok
}

func GetMeterfixMapLen() int {
return len(meterfix)
}

