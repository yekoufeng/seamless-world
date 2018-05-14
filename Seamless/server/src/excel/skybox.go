package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type SkyboxData struct {
Id uint64
}

var skybox map[uint64]SkyboxData
var skyboxLock sync.RWMutex

func LoadSkybox(){
skyboxLock.Lock()
defer skyboxLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/skybox.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &skybox)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetSkyboxMap() map[uint64]SkyboxData {
skyboxLock.RLock()
defer skyboxLock.RUnlock()

skybox2 := make(map[uint64]SkyboxData)
for k, v := range skybox{
skybox2[k] = v
}

return skybox2
}

func GetSkybox(key uint64) (SkyboxData, bool) {
skyboxLock.RLock()
defer skyboxLock.RUnlock()

val, ok := skybox[key]

return val, ok
}

func GetSkyboxMapLen() int {
return len(skybox)
}

