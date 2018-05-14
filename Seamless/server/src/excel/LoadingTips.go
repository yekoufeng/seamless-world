package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type LoadingTipsData struct {
Id uint64
Content string
}

var LoadingTips map[uint64]LoadingTipsData
var LoadingTipsLock sync.RWMutex

func LoadLoadingTips(){
LoadingTipsLock.Lock()
defer LoadingTipsLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/LoadingTips.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &LoadingTips)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetLoadingTipsMap() map[uint64]LoadingTipsData {
LoadingTipsLock.RLock()
defer LoadingTipsLock.RUnlock()

LoadingTips2 := make(map[uint64]LoadingTipsData)
for k, v := range LoadingTips{
LoadingTips2[k] = v
}

return LoadingTips2
}

func GetLoadingTips(key uint64) (LoadingTipsData, bool) {
LoadingTipsLock.RLock()
defer LoadingTipsLock.RUnlock()

val, ok := LoadingTips[key]

return val, ok
}

func GetLoadingTipsMapLen() int {
return len(LoadingTips)
}

