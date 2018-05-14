package main

import (
	"sync"
	"zeus/common"

	log "github.com/cihub/seelog"
)

type User struct {
	IsMainRole bool
	EntityID   uint64
	EntityType string

	Props
	UserState

	UserView
}

func NewUser() *User {
	u := &User{
		Props: Props{
			propsMutex: &sync.Mutex{},
		},
		UserState: UserState{
			speed: 60.0,
			mtx:   &sync.Mutex{},
		},
	}

	u.um = NewUserModel(GetMainView().window.Render)

	GetMouseableMgr().AddMouseable(&u.UserView)

	return u
}

func (user *User) ParseProps(num uint16, data []byte) {
	stream := common.NewByteStream(data)
	for i := uint16(0); i < num; i++ {
		name, err := stream.ReadStr()
		if err != nil {
			log.Error(err)
			return
		}
		user.ReadValueFromStream(name, stream)
	}
}

func (user *User) Loop() {
	if user.IsMainRole {
		user.UpdateMove()
		user.SyncMoveToCell()
	}
}

func (user *User) Destroy() {
	GetMouseableMgr().RemoveMouseable(&user.UserView)
}
