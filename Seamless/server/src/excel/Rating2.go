package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type Rating2Data struct {
Id uint64
InitWinRating uint64
InitKillRating uint64
InitFinalRating uint64
Range1Min uint64
Range1Max uint64
Range2Min uint64
Range2Max uint64
Range3Min uint64
NValueOfRange1 float32
MValueOfRange1 float32
CValueOfRange1 float32
RiValueOfRange1 float32
NValueOfRange2 float32
MValueOfRange2 float32
CValueOfRange2 float32
RiValueOfRange2 float32
NValueOfRange3 float32
MValueOfRange3 float32
CValueOfRange3 float32
RiValueOfRange3 float32
KillKValue float32
WinRatingPer float32
KillRatingPer float32
HValueOfCorrection string
DValue float32
LValue float32
IncreaseTValue int
}

var Rating2 map[uint64]Rating2Data
var Rating2Lock sync.RWMutex

func LoadRating2(){
Rating2Lock.Lock()
defer Rating2Lock.Unlock()

data,err := ioutil.ReadFile("../res/excel/Rating2.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &Rating2)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetRating2Map() map[uint64]Rating2Data {
Rating2Lock.RLock()
defer Rating2Lock.RUnlock()

Rating22 := make(map[uint64]Rating2Data)
for k, v := range Rating2{
Rating22[k] = v
}

return Rating22
}

func GetRating2(key uint64) (Rating2Data, bool) {
Rating2Lock.RLock()
defer Rating2Lock.RUnlock()

val, ok := Rating2[key]

return val, ok
}

func GetRating2MapLen() int {
return len(Rating2)
}

