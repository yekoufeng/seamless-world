package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type ResultcoinData struct {
Id uint64
Single uint64
Two uint64
Four uint64
}

var resultcoin map[uint64]ResultcoinData
var resultcoinLock sync.RWMutex

func LoadResultcoin(){
resultcoinLock.Lock()
defer resultcoinLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/resultcoin.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &resultcoin)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetResultcoinMap() map[uint64]ResultcoinData {
resultcoinLock.RLock()
defer resultcoinLock.RUnlock()

resultcoin2 := make(map[uint64]ResultcoinData)
for k, v := range resultcoin{
resultcoin2[k] = v
}

return resultcoin2
}

func GetResultcoin(key uint64) (ResultcoinData, bool) {
resultcoinLock.RLock()
defer resultcoinLock.RUnlock()

val, ok := resultcoin[key]

return val, ok
}

func GetResultcoinMapLen() int {
return len(resultcoin)
}

