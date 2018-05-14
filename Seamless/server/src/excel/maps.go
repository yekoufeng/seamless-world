package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type MapsData struct {
Id uint64
Name string
Born_pos string
Water_height float32
Saferadius float32
Safe_zone string
Refresh_item string
Width float32
Height float32
Parachute_height float32
Fly_Speed uint64
Parts_num uint64
End_point uint64
E_start uint64
E_end uint64
W_start uint64
W_end uint64
N_start uint64
N_end uint64
S_start uint64
S_end uint64
Drop_Time_Point float32
AI_Drop_Delay uint64
AI_Drop_Percent float32
AI_Drop_Time float32
Ai_Time uint64
Fog_height uint64
Bomb_zone string
Skybox string
}

var maps map[uint64]MapsData
var mapsLock sync.RWMutex

func LoadMaps(){
mapsLock.Lock()
defer mapsLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/maps.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &maps)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetMapsMap() map[uint64]MapsData {
mapsLock.RLock()
defer mapsLock.RUnlock()

maps2 := make(map[uint64]MapsData)
for k, v := range maps{
maps2[k] = v
}

return maps2
}

func GetMaps(key uint64) (MapsData, bool) {
mapsLock.RLock()
defer mapsLock.RUnlock()

val, ok := maps[key]

return val, ok
}

func GetMapsMapLen() int {
return len(maps)
}

