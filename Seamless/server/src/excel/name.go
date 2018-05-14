package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type NameData struct {
Nameid uint64
Name string
}

var name map[uint64]NameData
var nameLock sync.RWMutex

func LoadName(){
nameLock.Lock()
defer nameLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/name.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &name)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetNameMap() map[uint64]NameData {
nameLock.RLock()
defer nameLock.RUnlock()

name2 := make(map[uint64]NameData)
for k, v := range name{
name2[k] = v
}

return name2
}

func GetName(key uint64) (NameData, bool) {
nameLock.RLock()
defer nameLock.RUnlock()

val, ok := name[key]

return val, ok
}

func GetNameMapLen() int {
return len(name)
}

