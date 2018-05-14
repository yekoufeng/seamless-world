package main

import (
	//"zeus/env"
	"zeus/zlog"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()

	zlog.Init()

	if !env.Load("../res/config/server.json") {
		log.Error("加载配置文件失败")
		return
	}

	GetSimpleSrvInst().Run()
}
