package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type CharacterstateData struct {
Id uint64
Name string
BeDetectedFix float32
BeHitFix float32
AiViewFix float32
}

var characterstate map[uint64]CharacterstateData
var characterstateLock sync.RWMutex

func LoadCharacterstate(){
characterstateLock.Lock()
defer characterstateLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/characterstate.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &characterstate)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetCharacterstateMap() map[uint64]CharacterstateData {
characterstateLock.RLock()
defer characterstateLock.RUnlock()

characterstate2 := make(map[uint64]CharacterstateData)
for k, v := range characterstate{
characterstate2[k] = v
}

return characterstate2
}

func GetCharacterstate(key uint64) (CharacterstateData, bool) {
characterstateLock.RLock()
defer characterstateLock.RUnlock()

val, ok := characterstate[key]

return val, ok
}

func GetCharacterstateMapLen() int {
return len(characterstate)
}

