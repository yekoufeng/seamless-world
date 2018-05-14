package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type MapitemData struct {
Id uint64
Zone string
Rate uint64
}

var mapitem map[uint64]MapitemData
var mapitemLock sync.RWMutex

func LoadMapitem(){
mapitemLock.Lock()
defer mapitemLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/mapitem.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &mapitem)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetMapitemMap() map[uint64]MapitemData {
mapitemLock.RLock()
defer mapitemLock.RUnlock()

mapitem2 := make(map[uint64]MapitemData)
for k, v := range mapitem{
mapitem2[k] = v
}

return mapitem2
}

func GetMapitem(key uint64) (MapitemData, bool) {
mapitemLock.RLock()
defer mapitemLock.RUnlock()

val, ok := mapitem[key]

return val, ok
}

func GetMapitemMapLen() int {
return len(mapitem)
}

