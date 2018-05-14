package excel

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
)

type LoadFunc func()

var loadFuncs = map[string]LoadFunc{
"AIbirthdeath_SafeCircle":LoadAIbirthdeath_SafeCircle,
"AIbirthdeath_Server":LoadAIbirthdeath_Server,
"AImatch":LoadAImatch,
	"Sadscore":       LoadSadscore,
	"AInumber":       LoadAInumber,
	"LoadingTips":    LoadLoadingTips,
	"Errorcode":      LoadErrorcode,
	"Item":           LoadItem,
	"Mapitem":        LoadMapitem,
	"Mapitemrate":    LoadMapitemrate,
	"Logic":          LoadLogic,
	"Carrier":        LoadCarrier,
	"Bombrule":       LoadBombrule,
	"Parachuting":    LoadParachuting,
	"Meterfix":       LoadMeterfix,
	"Season":         LoadSeason,
	"Rating":         LoadRating,
	"Role":           LoadRole,
	"Energy":         LoadEnergy,
	"Adkey":          LoadAdkey,
	"Barrier":        LoadBarrier,
	"Gun":            LoadGun,
	"Clothes":        LoadClothes,
	"Phonenumber":    LoadPhonenumber,
	"Achievement":    LoadAchievement,
	"Detail":         LoadDetail,
	"Keywords":       LoadKeywords,
	"Characterstate": LoadCharacterstate,
	"Skybox":         LoadSkybox,
	"Maps":           LoadMaps,
	"Maprule":        LoadMaprule,
	"Store":          LoadStore,
	"Resultcoin":     LoadResultcoin,
	"Name":           LoadName,
	"Birthplace":     LoadBirthplace,
	"Pay":            LoadPay,
	"Rating2":        LoadRating2,
	"System":         LoadSystem,
	"Ai":             LoadAi,
}

func init() {
	for _, f := range loadFuncs {
		f()
	}

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error(err)
			return
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case ev := <-watcher.Events:
					log.Infof("Watch %s Op %s", ev.Name, ev.Op)

					if strings.Contains(ev.Name, "reload.json") {
						if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create {
							log.Info("Start Reload!")

							if err := Reload("../res/excel/reload.json"); err != nil {
								log.Error("Reload failed ", err)
							} else {
								log.Info("Reload excel success")
							}
						}
					}
				case err := <-watcher.Errors:
					log.Error(err)
				}
			}
		}()

		if err = watcher.Add("../res/excel/"); err != nil {
			log.Error(err)
		} else {
			log.Info("Watch excel path")
		}

		<-done
	}()

}

func Reload(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config := make(map[string]bool)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	for k, v := range config {
		if v {
			f, ok := loadFuncs[k]
			if ok {
				log.Info("Reload ", k)
				f()
			} else {
				log.Error("Reload ", k, " failed, can't find the load function")
			}
		}
	}

	return nil
}

func ReloadAll() {
	for _, f := range loadFuncs {
		f()
	}
}
