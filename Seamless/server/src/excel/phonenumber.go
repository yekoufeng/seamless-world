package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type PhonenumberData struct {
Id uint64
Phonenumber string
}

var phonenumber map[uint64]PhonenumberData
var phonenumberLock sync.RWMutex

func LoadPhonenumber(){
phonenumberLock.Lock()
defer phonenumberLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/phonenumber.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &phonenumber)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetPhonenumberMap() map[uint64]PhonenumberData {
phonenumberLock.RLock()
defer phonenumberLock.RUnlock()

phonenumber2 := make(map[uint64]PhonenumberData)
for k, v := range phonenumber{
phonenumber2[k] = v
}

return phonenumber2
}

func GetPhonenumber(key uint64) (PhonenumberData, bool) {
phonenumberLock.RLock()
defer phonenumberLock.RUnlock()

val, ok := phonenumber[key]

return val, ok
}

func GetPhonenumberMapLen() int {
return len(phonenumber)
}

