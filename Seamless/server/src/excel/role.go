package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type RoleData struct {
Id uint64
Name string
Initmvspeed float32
Mvspeed string
SlowMvSpeed string
Crawlspeed string
SlowCrawlSpeed string
Initswimspeed float32
Swimspeed string
SlowSwimSpeed string
WillDieSpeed string
CrouchSpeed string
SlowCrouchSpeed string
FastRunSpeed float32
Inithp float32
Initattack int
Gender int
Mvspeedlimit float32
Crawlspeedlimit float32
Swimspeedlimit float32
Crouchspeedlimit float32
Vehiclespeedlimit float32
Willdiespeedlimit float32
Fallorjumpspeedlimit float32
ChangeClothes int
HeadSkin int
FaceSkin int
TopsSkin int
PantSkin int
ShoesSkin int
}

var role map[uint64]RoleData
var roleLock sync.RWMutex

func LoadRole(){
roleLock.Lock()
defer roleLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/role.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &role)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetRoleMap() map[uint64]RoleData {
roleLock.RLock()
defer roleLock.RUnlock()

role2 := make(map[uint64]RoleData)
for k, v := range role{
role2[k] = v
}

return role2
}

func GetRole(key uint64) (RoleData, bool) {
roleLock.RLock()
defer roleLock.RUnlock()

val, ok := role[key]

return val, ok
}

func GetRoleMapLen() int {
return len(role)
}

