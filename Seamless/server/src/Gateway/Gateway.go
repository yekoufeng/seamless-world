package main

import (
	"net/http"
	_ "net/http/pprof"
	"zeus/zlog"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("../res/config/server.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}

	logDir := viper.GetString("Config.LogDir")
	logLevel := viper.GetString("Config.LogLevel")
	zlog.Init(logDir, logLevel)
	defer log.Flush()

	if viper.GetBool("Config.Debug") {
		go func() {
			http.ListenAndServe("localhost:6051", nil)
		}()
	}

	GetSrvInst().Run()
}
