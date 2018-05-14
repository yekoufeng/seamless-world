package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type CarrierData struct {
Id uint64
Name string
Seat uint64
FuelCap float32
FuelConsume float32
FuelMin float32
FuelMax float32
}

var carrier map[uint64]CarrierData
var carrierLock sync.RWMutex

func LoadCarrier(){
carrierLock.Lock()
defer carrierLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/carrier.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &carrier)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetCarrierMap() map[uint64]CarrierData {
carrierLock.RLock()
defer carrierLock.RUnlock()

carrier2 := make(map[uint64]CarrierData)
for k, v := range carrier{
carrier2[k] = v
}

return carrier2
}

func GetCarrier(key uint64) (CarrierData, bool) {
carrierLock.RLock()
defer carrierLock.RUnlock()

val, ok := carrier[key]

return val, ok
}

func GetCarrierMapLen() int {
return len(carrier)
}

