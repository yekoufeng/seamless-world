package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type AchievementData struct {
Id uint64
Type uint64
Name string
Content string
BonusID uint64
BonusNum uint64
MailTitle string
MailContent string
}

var achievement map[uint64]AchievementData
var achievementLock sync.RWMutex

func LoadAchievement(){
achievementLock.Lock()
defer achievementLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/achievement.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &achievement)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetAchievementMap() map[uint64]AchievementData {
achievementLock.RLock()
defer achievementLock.RUnlock()

achievement2 := make(map[uint64]AchievementData)
for k, v := range achievement{
achievement2[k] = v
}

return achievement2
}

func GetAchievement(key uint64) (AchievementData, bool) {
achievementLock.RLock()
defer achievementLock.RUnlock()

val, ok := achievement[key]

return val, ok
}

func GetAchievementMapLen() int {
return len(achievement)
}

