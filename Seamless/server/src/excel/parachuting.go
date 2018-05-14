package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type ParachutingData struct {
Id uint64
Sidelen int
Notallowjump string
Row int
Column int
Len int
BeforeGameTime int
}

var parachuting map[uint64]ParachutingData
var parachutingLock sync.RWMutex

func LoadParachuting(){
parachutingLock.Lock()
defer parachutingLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/parachuting.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &parachuting)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetParachutingMap() map[uint64]ParachutingData {
parachutingLock.RLock()
defer parachutingLock.RUnlock()

parachuting2 := make(map[uint64]ParachutingData)
for k, v := range parachuting{
parachuting2[k] = v
}

return parachuting2
}

func GetParachuting(key uint64) (ParachutingData, bool) {
parachutingLock.RLock()
defer parachutingLock.RUnlock()

val, ok := parachuting[key]

return val, ok
}

func GetParachutingMapLen() int {
return len(parachuting)
}

