package main

import (
	"zeus/iserver"

	log "github.com/cihub/seelog"
)

var usermgr *GateUserMgr

// GetUserMgr 获取网关玩家管理器
func GetUserMgr() *GateUserMgr {
	if usermgr == nil {
		usermgr = &GateUserMgr{}
	}

	return usermgr
}

// GateUserMgr GateUser 管理
type GateUserMgr struct {
}

func (mgr *GateUserMgr) getUserByUID(uid uint64) *GateUser {
	if user, ok := GetSrvInst().GetEntityByDBID("Player", uid).(*GateUser); ok {
		return user
	}
	return nil
}

func (mgr *GateUserMgr) getUserByID(id uint64) *GateUser {
	if user, ok := GetSrvInst().GetEntity(id).(*GateUser); ok {
		return user
	}
	return nil
}

func (mgr *GateUserMgr) removeUserByUID(uid uint64) error {
	return GetSrvInst().DestroyEntityByDBID("Player", uid)
}

func (mgr *GateUserMgr) addUser(uid uint64) *GateUser {

	user, err := GetSrvInst().CreateEntityAll("Player", uid, "", true)
	if err != nil {
		log.Error("Add user failed", err, "UID:", uid)
		return nil
	}

	return user.(*GateUser)
}

func (mgr *GateUserMgr) forceLogout() {

	users := make([]*GateUser, 0, 100)

	//防止迭代器失效
	GetSrvInst().TravsalEntity("Player", func(e iserver.IEntity) {

		user := e.(*GateUser)
		if user == nil {
			return
		}

		users = append(users, user)
	})

	for _, u := range users {
		u.logout()
	}
}
