package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
)

// GmMgr Gm命令
type GmMgr struct {
	user *LobbyUser

	// cmds 命令集合
	cmds map[string](func(map[string]string))

	gmSkyType      uint32 // gm设置天空盒类型
	isUseGmSkyType uint32 // 0 不使用 1 使用
}

// NewGmMgr 获取Gm管理器
func NewGmMgr(user *LobbyUser) *GmMgr {
	gm := &GmMgr{
		user:           user,
		cmds:           make(map[string](func(map[string]string))),
		gmSkyType:      0,
		isUseGmSkyType: 0,
	}
	gm.init()

	return gm
}

// init 初始化管理器
func (gm *GmMgr) init() {
	gm.cmds["SetSkyType"] = gm.SetSkyType // 设置天空盒类型
}

// exec 执行命令
func (gm *GmMgr) exec(paras string) {
	pairSet := strings.Split(paras, " ")

	pairMap := make(map[string]string)
	for _, pair := range pairSet {
		paraSet := strings.Split(pair, "=")
		if len(paraSet) != 2 {
			continue
		}

		pairMap[paraSet[0]] = paraSet[1]
	}

	if cmdStr, ok := pairMap["lcmd"]; ok {
		if cmd, ok := gm.cmds[cmdStr]; ok {
			cmd(pairMap)
		} else {
			fmt.Println("命令集cmds不包含命令", cmdStr)
		}
	} else {
		fmt.Println("前端参数不包含命令!")
	}
}

// SetSkyType 设置天空盒类型
func (gm *GmMgr) SetSkyType(paras map[string]string) {
	log.Info("SetSkyType")

	var skyType, useGm int
	if strSkyType, typeErr := paras["skyType"]; typeErr {
		var valueErr error
		skyType, valueErr = strconv.Atoi(strSkyType)
		if valueErr != nil {
			log.Error(valueErr)
			return
		}
	} else {
		log.Error(typeErr)
		return
	}

	if strIsUseType, isUseErr := paras["isUseGm"]; isUseErr {
		var valueErr error
		useGm, valueErr = strconv.Atoi(strIsUseType)
		if valueErr != nil {
			log.Error(valueErr)
			return
		}
	} else {
		log.Error(isUseErr)
		return
	}

	gm.gmSkyType = uint32(skyType)
	gm.isUseGmSkyType = uint32(useGm)

	log.Error("SetSpaceType :%+v, %d, %d", paras, uint32(skyType), uint32(useGm))
}
