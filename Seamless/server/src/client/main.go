package main

import (
	"common"
	"flag"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"zeus/httppprof"
	"zeus/zlog"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

var closeChan = make(chan struct{})
var exit = int32(0)

func ExitClient() {
	if atomic.CompareAndSwapInt32(&exit, 0, 1) {
		close(closeChan)
	}
}

// var viewStart = make(chan struct{})

func main() {

	flag.Parse()

	pprofPort := "11181"
	if *httpport != "" {
		pprofPort = *httpport
	}
	go httppprof.StartPProf(":" + pprofPort)

	viper.SetConfigFile("../res/config/server.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}

	zlog.Init("../log/", "debug")
	defer log.Flush()

	GetAssetsMgr().PreLoad(GetMainView().window.Render)
	GetMainView().Init()

	common.InitMsg()

	if err := GetClient().Login(); err != nil {
		log.Error(err)
		return
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		<-c
		log.Debug("Client Exist by signal")
		ExitClient()
	}()

	go GetClient().Loop()

	// <-viewStart
	GetMainView().Start()

	log.Debug("Client Exist")
}
