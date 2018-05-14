package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AdkeyData struct {
Id uint64
Adkey string
}

var adkey map[uint64]AdkeyData
var adkeyLock sync.RWMutex

func LoadAdkey(){
adkeyLock.Lock()
defer adkeyLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/adkey.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &adkey)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAdkeyMap() map[uint64]AdkeyData {
adkeyLock.RLock()
defer adkeyLock.RUnlock()

adkey2 := make(map[uint64]AdkeyData)
for k, v := range adkey{
adkey2[k] = v
}

return adkey2
}

func GetAdkey(key uint64) (AdkeyData, bool) {
adkeyLock.RLock()
defer adkeyLock.RUnlock()

val, ok := adkey[key]

return val, ok
}

func GetAdkeyMapLen() int {
return len(adkey)
}

