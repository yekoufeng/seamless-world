package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type DetailData struct {
Id uint64
Name string
Dailylogin string
EndPlus string
DailyLoginMailTitle string
DailyLoginMail string
}

var detail map[uint64]DetailData
var detailLock sync.RWMutex

func LoadDetail(){
detailLock.Lock()
defer detailLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/detail.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &detail)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetDetailMap() map[uint64]DetailData {
detailLock.RLock()
defer detailLock.RUnlock()

detail2 := make(map[uint64]DetailData)
for k, v := range detail{
detail2[k] = v
}

return detail2
}

func GetDetail(key uint64) (DetailData, bool) {
detailLock.RLock()
defer detailLock.RUnlock()

val, ok := detail[key]

return val, ok
}

func GetDetailMapLen() int {
return len(detail)
}

