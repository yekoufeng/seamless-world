package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type ClothesData struct {
Id uint64
Name string
Position int
Gender int
DropItemId int
}

var clothes map[uint64]ClothesData
var clothesLock sync.RWMutex

func LoadClothes(){
clothesLock.Lock()
defer clothesLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/clothes.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &clothes)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetClothesMap() map[uint64]ClothesData {
clothesLock.RLock()
defer clothesLock.RUnlock()

clothes2 := make(map[uint64]ClothesData)
for k, v := range clothes{
clothes2[k] = v
}

return clothes2
}

func GetClothes(key uint64) (ClothesData, bool) {
clothesLock.RLock()
defer clothesLock.RUnlock()

val, ok := clothes[key]

return val, ok
}

func GetClothesMapLen() int {
return len(clothes)
}

