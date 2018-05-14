package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type SadscoreData struct {
Id uint64
Ranking int
Sadscore float32
}

var sadscore map[uint64]SadscoreData
var sadscoreLock sync.RWMutex

func LoadSadscore(){
sadscoreLock.Lock()
defer sadscoreLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/sadscore.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &sadscore)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetSadscoreMap() map[uint64]SadscoreData {
sadscoreLock.RLock()
defer sadscoreLock.RUnlock()

sadscore2 := make(map[uint64]SadscoreData)
for k, v := range sadscore{
sadscore2[k] = v
}

return sadscore2
}

func GetSadscore(key uint64) (SadscoreData, bool) {
sadscoreLock.RLock()
defer sadscoreLock.RUnlock()

val, ok := sadscore[key]

return val, ok
}

func GetSadscoreMapLen() int {
return len(sadscore)
}

