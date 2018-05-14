package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AIequipmentData struct {
SafeCircle int
Helmet_0 int
Helmet_1202 int
Helmet_1204 int
Helmet_1206 int
Armor_0 int
Armor_1201 int
Armor_1203 int
Armor_1205 int
Melee_0 int
Melee_10001 int
Melee_10002 int
Melee_10003 int
Melee_10004 int
Gun_kind_10200_10100 int
Gun_kind_10300_10100 int
Gun_kind_10100_10400 int
Gun_kind_10500_10100 int
Gun_kind_10300_10200 int
Gun_kind_10200_10400 int
Gun_kind_10200_10500 int
Gun_kind_10300_10400 int
Gun_kind_10300_10500 int
Gun_kind_10500_10400 int
BulletNum_10100 string
BulletNum_10200 string
BulletNum_10300 string
BulletNum_10400 string
BulletNum_10500 string
ReformItemsNum_0 int
ReformItemsNum_1 int
ReformItemsNum_2 int
ReformItemsNum_3 int
ReformItemsNum_4 int
ReformItemsNum_5 int
ReformItemsNum_6 int
GrenadeNum_0 int
GrenadeNum_1 int
GrenadeNum_2 int
GrenadeNum_3 int
GrenadeNum_4 int
GrenadeNum_5 int
GrenadeNum_6 int
GrenadeNum_7 int
GrenadeNum_8 int
GrenadeNum_9 int
GrenadeNum_10 int
Grenade_2100 int
Grenade_2101 int
Grenade_2102 int
MedicinesNum_0 int
MedicinesNum_1 int
MedicinesNum_2 int
MedicinesNum_3 int
MedicinesNum_4 int
MedicinesNum_5 int
MedicinesNum_6 int
MedicinesNum_7 int
MedicinesNum_8 int
MedicinesNum_9 int
MedicinesNum_10 int
Medicine_1101 int
Medicine_1102 int
Medicine_1117 int
Medicine_1118 int
Medicine_1119 int
Medicine_1120 int
ClothesNum_0 int
ClothesNum_1 int
ClothesNum_2 int
ClothesNum_3 int
ClothesNum_4 int
ClothesNum_5 int
Clothes_3000 int
Clothes_3500 int
Clothes_4000 int
Clothes_4500 int
Clothes_5000 int
}

var AIequipment map[uint64]AIequipmentData
var AIequipmentLock sync.RWMutex

func LoadAIequipment(){
AIequipmentLock.Lock()
defer AIequipmentLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/AIequipment.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &AIequipment)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAIequipmentMap() map[uint64]AIequipmentData {
AIequipmentLock.RLock()
defer AIequipmentLock.RUnlock()

AIequipment2 := make(map[uint64]AIequipmentData)
for k, v := range AIequipment{
AIequipment2[k] = v
}

return AIequipment2
}

func GetAIequipment(key uint64) (AIequipmentData, bool) {
AIequipmentLock.RLock()
defer AIequipmentLock.RUnlock()

val, ok := AIequipment[key]

return val, ok
}

func GetAIequipmentMapLen() int {
return len(AIequipment)
}

