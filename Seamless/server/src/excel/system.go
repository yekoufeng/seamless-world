package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type SystemData struct {
Id uint64
Value uint64
Name string
}

var system map[uint64]SystemData
var systemLock sync.RWMutex

func LoadSystem(){
systemLock.Lock()
defer systemLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/system.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &system)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetSystemMap() map[uint64]SystemData {
systemLock.RLock()
defer systemLock.RUnlock()

system2 := make(map[uint64]SystemData)
for k, v := range system{
system2[k] = v
}

return system2
}

func GetSystem(key uint64) (SystemData, bool) {
systemLock.RLock()
defer systemLock.RUnlock()

val, ok := system[key]

return val, ok
}

func GetSystemMapLen() int {
return len(system)
}

