package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AIbirthdeath_SafeCircleData struct {
SafeCircle int
SummonAITime_Min int
SummonAITime_Max int
PlayerControlAINum int
ExpectPlayerNum int
MinPlayerNum int
AIDeath_OutOfTime int
AIDeath_Traffic int
AIDeath_Melee int
AIDeath_Shoot int
AIDeath_Grenade int
AIDeath_Fall int
AIDeath_Traffic_VehicleStrick int
AIDeath_Traffic_JumpVehicle int
AIDeath_Traffic_VehicleExplode int
AIDeath_Traffic_VehicleBurn int
AIDeath_Shoot_ShotGun int
AIDeath_Shoot_Pistol int
AIDeath_Shoot_SubmachineGun int
AIDeath_Shoot_AssaultRifle int
AIDeath_Shoot_SniperRifle int
}

var AIbirthdeath_SafeCircle map[uint64]AIbirthdeath_SafeCircleData
var AIbirthdeath_SafeCircleLock sync.RWMutex

func LoadAIbirthdeath_SafeCircle(){
AIbirthdeath_SafeCircleLock.Lock()
defer AIbirthdeath_SafeCircleLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/AIbirthdeath_SafeCircle.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &AIbirthdeath_SafeCircle)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAIbirthdeath_SafeCircleMap() map[uint64]AIbirthdeath_SafeCircleData {
AIbirthdeath_SafeCircleLock.RLock()
defer AIbirthdeath_SafeCircleLock.RUnlock()

AIbirthdeath_SafeCircle2 := make(map[uint64]AIbirthdeath_SafeCircleData)
for k, v := range AIbirthdeath_SafeCircle{
AIbirthdeath_SafeCircle2[k] = v
}

return AIbirthdeath_SafeCircle2
}

func GetAIbirthdeath_SafeCircle(key uint64) (AIbirthdeath_SafeCircleData, bool) {
AIbirthdeath_SafeCircleLock.RLock()
defer AIbirthdeath_SafeCircleLock.RUnlock()

val, ok := AIbirthdeath_SafeCircle[key]

return val, ok
}

func GetAIbirthdeath_SafeCircleMapLen() int {
return len(AIbirthdeath_SafeCircle)
}

