package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type PayData struct {
Id uint64
Productid string
Goodslocalicon string
Buynum string
Payitem string
Ext string
IsChange uint64
Sername string
ServiceCode string
AutoPay uint64
ResId string
MallLogo string
PayType uint64
}

var pay map[uint64]PayData
var payLock sync.RWMutex

func LoadPay(){
payLock.Lock()
defer payLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/pay.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &pay)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetPayMap() map[uint64]PayData {
payLock.RLock()
defer payLock.RUnlock()

pay2 := make(map[uint64]PayData)
for k, v := range pay{
pay2[k] = v
}

return pay2
}

func GetPay(key uint64) (PayData, bool) {
payLock.RLock()
defer payLock.RUnlock()

val, ok := pay[key]

return val, ok
}

func GetPayMapLen() int {
return len(pay)
}

