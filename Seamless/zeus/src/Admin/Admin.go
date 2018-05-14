package main

import (
	"flag"
	//"zeus/env"
	"zeus/zlog"

	log "github.com/cihub/seelog"
)

func main() {
	flag.Parse()
	zlog.Init()
	defer log.Flush()

	if !env.Load("../res/config/server.json") {
		log.Error("加载配置文件失败")
		return
	}

	GetAdminServer().Init()
	GetAdminServer().Start()
}
