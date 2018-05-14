package excel

import (
"encoding/json"
"io/ioutil"
"sync"

log "github.com/cihub/seelog"
)
type KeywordsData struct {
Nameid uint64
Name string
}

var keywords map[uint64]KeywordsData
var keywordsLock sync.RWMutex

func LoadKeywords(){
keywordsLock.Lock()
defer keywordsLock.Unlock()

data,err := ioutil.ReadFile("../res/excel/keywords.json")
if err != nil {
log.Error("ReadFile err: ", err)
return
}

err = json.Unmarshal(data, &keywords)
if err != nil {
log.Error("Unmarshal err: ", err)
return
}
}

func GetKeywordsMap() map[uint64]KeywordsData {
keywordsLock.RLock()
defer keywordsLock.RUnlock()

keywords2 := make(map[uint64]KeywordsData)
for k, v := range keywords{
keywords2[k] = v
}

return keywords2
}

func GetKeywords(key uint64) (KeywordsData, bool) {
keywordsLock.RLock()
defer keywordsLock.RUnlock()

val, ok := keywords[key]

return val, ok
}

func GetKeywordsMapLen() int {
return len(keywords)
}

