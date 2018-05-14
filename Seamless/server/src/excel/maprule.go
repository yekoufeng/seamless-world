package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type MapruleData struct {
Id uint64
Diameter uint64
Length uint64
Ifland uint64
Shrinkarea uint64
Shrinktime uint64
Shrinkdam uint64
Boxrefresh uint64
Bombrefresh uint64
Delaybombdam uint64
Bombtime uint64
Bombradius float32
Bombdam uint64
Bombspeed float32
Bombdamrradius float32
Aidownmove float32
}

var maprule map[uint64]MapruleData
var mapruleLock sync.RWMutex

func LoadMaprule(){
mapruleLock.Lock()
defer mapruleLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/maprule.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &maprule)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetMapruleMap() map[uint64]MapruleData {
mapruleLock.RLock()
defer mapruleLock.RUnlock()

maprule2 := make(map[uint64]MapruleData)
for k, v := range maprule{
maprule2[k] = v
}

return maprule2
}

func GetMaprule(key uint64) (MapruleData, bool) {
mapruleLock.RLock()
defer mapruleLock.RUnlock()

val, ok := maprule[key]

return val, ok
}

func GetMapruleMapLen() int {
return len(maprule)
}

