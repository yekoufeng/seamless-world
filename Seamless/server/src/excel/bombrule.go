package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type BombruleData struct {
Id uint64
Bombarea uint64
Bombtime uint64
Bombradius float32
Bombstart uint64
Bombspeed float32
Bombdam uint64
Bombdamrradius float32
}

var bombrule map[uint64]BombruleData
var bombruleLock sync.RWMutex

func LoadBombrule(){
bombruleLock.Lock()
defer bombruleLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/bombrule.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &bombrule)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetBombruleMap() map[uint64]BombruleData {
bombruleLock.RLock()
defer bombruleLock.RUnlock()

bombrule2 := make(map[uint64]BombruleData)
for k, v := range bombrule{
bombrule2[k] = v
}

return bombrule2
}

func GetBombrule(key uint64) (BombruleData, bool) {
bombruleLock.RLock()
defer bombruleLock.RUnlock()

val, ok := bombrule[key]

return val, ok
}

func GetBombruleMapLen() int {
return len(bombrule)
}

