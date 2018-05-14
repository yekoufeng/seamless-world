package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type RatingData struct {
Id uint64
Value float32
Name string
}

var rating map[uint64]RatingData
var ratingLock sync.RWMutex

func LoadRating(){
ratingLock.Lock()
defer ratingLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/rating.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &rating)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetRatingMap() map[uint64]RatingData {
ratingLock.RLock()
defer ratingLock.RUnlock()

rating2 := make(map[uint64]RatingData)
for k, v := range rating{
rating2[k] = v
}

return rating2
}

func GetRating(key uint64) (RatingData, bool) {
ratingLock.RLock()
defer ratingLock.RUnlock()

val, ok := rating[key]

return val, ok
}

func GetRatingMapLen() int {
return len(rating)
}

