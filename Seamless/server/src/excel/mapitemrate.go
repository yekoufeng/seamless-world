package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type MapitemrateData struct {
Id uint64
Itemlist1 string
Itemlist2 string
Itemlist3 string
Itemlist4 string
Itemlist5 string
Itemlist6 string
Itemlist7 string
Itemlist8 string
Itemlist9 string
Itemlist10 string
}

var mapitemrate map[uint64]MapitemrateData
var mapitemrateLock sync.RWMutex

func LoadMapitemrate(){
mapitemrateLock.Lock()
defer mapitemrateLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/mapitemrate.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &mapitemrate)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetMapitemrateMap() map[uint64]MapitemrateData {
mapitemrateLock.RLock()
defer mapitemrateLock.RUnlock()

mapitemrate2 := make(map[uint64]MapitemrateData)
for k, v := range mapitemrate{
mapitemrate2[k] = v
}

return mapitemrate2
}

func GetMapitemrate(key uint64) (MapitemrateData, bool) {
mapitemrateLock.RLock()
defer mapitemrateLock.RUnlock()

val, ok := mapitemrate[key]

return val, ok
}

func GetMapitemrateMapLen() int {
return len(mapitemrate)
}

