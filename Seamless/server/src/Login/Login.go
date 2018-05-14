package main

import (
	"os"
	"os/signal"
	"syscall"
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

	srv := NewServer()
	srv.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-c
	srv.Stop()
}
