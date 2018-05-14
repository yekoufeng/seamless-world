package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AIbirthdeath_ServerData struct {
Id uint64
Name string
Value float32
}

var AIbirthdeath_Server map[uint64]AIbirthdeath_ServerData
var AIbirthdeath_ServerLock sync.RWMutex

func LoadAIbirthdeath_Server(){
AIbirthdeath_ServerLock.Lock()
defer AIbirthdeath_ServerLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/AIbirthdeath_Server.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &AIbirthdeath_Server)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAIbirthdeath_ServerMap() map[uint64]AIbirthdeath_ServerData {
AIbirthdeath_ServerLock.RLock()
defer AIbirthdeath_ServerLock.RUnlock()

AIbirthdeath_Server2 := make(map[uint64]AIbirthdeath_ServerData)
for k, v := range AIbirthdeath_Server{
AIbirthdeath_Server2[k] = v
}

return AIbirthdeath_Server2
}

func GetAIbirthdeath_Server(key uint64) (AIbirthdeath_ServerData, bool) {
AIbirthdeath_ServerLock.RLock()
defer AIbirthdeath_ServerLock.RUnlock()

val, ok := AIbirthdeath_Server[key]

return val, ok
}

func GetAIbirthdeath_ServerMapLen() int {
return len(AIbirthdeath_Server)
}

