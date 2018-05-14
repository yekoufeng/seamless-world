package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type ErrorcodeData struct {
Id uint64
Name string
Content string
}

var errorcode map[uint64]ErrorcodeData
var errorcodeLock sync.RWMutex

func LoadErrorcode(){
errorcodeLock.Lock()
defer errorcodeLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/errorcode.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &errorcode)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetErrorcodeMap() map[uint64]ErrorcodeData {
errorcodeLock.RLock()
defer errorcodeLock.RUnlock()

errorcode2 := make(map[uint64]ErrorcodeData)
for k, v := range errorcode{
errorcode2[k] = v
}

return errorcode2
}

func GetErrorcode(key uint64) (ErrorcodeData, bool) {
errorcodeLock.RLock()
defer errorcodeLock.RUnlock()

val, ok := errorcode[key]

return val, ok
}

func GetErrorcodeMapLen() int {
return len(errorcode)
}

