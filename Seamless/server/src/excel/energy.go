package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type EnergyData struct {
Id uint64
Minvalue uint64
Value uint64
Addhp uint64
Addinterval uint64
Subenergy uint64
Subinterval uint64
Addspeed float32
}

var energy map[uint64]EnergyData
var energyLock sync.RWMutex

func LoadEnergy(){
energyLock.Lock()
defer energyLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/energy.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &energy)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetEnergyMap() map[uint64]EnergyData {
energyLock.RLock()
defer energyLock.RUnlock()

energy2 := make(map[uint64]EnergyData)
for k, v := range energy{
energy2[k] = v
}

return energy2
}

func GetEnergy(key uint64) (EnergyData, bool) {
energyLock.RLock()
defer energyLock.RUnlock()

val, ok := energy[key]

return val, ok
}

func GetEnergyMapLen() int {
return len(energy)
}

