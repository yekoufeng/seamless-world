package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type BirthplaceData struct {
Id uint64
X int
Z int
Length int
Width int
}

var birthplace map[uint64]BirthplaceData
var birthplaceLock sync.RWMutex

func LoadBirthplace(){
birthplaceLock.Lock()
defer birthplaceLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/birthplace.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &birthplace)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetBirthplaceMap() map[uint64]BirthplaceData {
birthplaceLock.RLock()
defer birthplaceLock.RUnlock()

birthplace2 := make(map[uint64]BirthplaceData)
for k, v := range birthplace{
birthplace2[k] = v
}

return birthplace2
}

func GetBirthplace(key uint64) (BirthplaceData, bool) {
birthplaceLock.RLock()
defer birthplaceLock.RUnlock()

val, ok := birthplace[key]

return val, ok
}

func GetBirthplaceMapLen() int {
return len(birthplace)
}

