package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type SeasonData struct {
Id uint64
Name string
EndTime string
}

var season map[uint64]SeasonData
var seasonLock sync.RWMutex

func LoadSeason(){
seasonLock.Lock()
defer seasonLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/season.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &season)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetSeasonMap() map[uint64]SeasonData {
seasonLock.RLock()
defer seasonLock.RUnlock()

season2 := make(map[uint64]SeasonData)
for k, v := range season{
season2[k] = v
}

return season2
}

func GetSeason(key uint64) (SeasonData, bool) {
seasonLock.RLock()
defer seasonLock.RUnlock()

val, ok := season[key]

return val, ok
}

func GetSeasonMapLen() int {
return len(season)
}

