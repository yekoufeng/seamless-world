package main

import (
	"sync"
	"zeus/msgdef"

	log "github.com/cihub/seelog"
)

type AOIS struct {
	users map[uint64]*User
	mtx   *sync.Mutex
}

func NewAOIS() *AOIS {
	return &AOIS{
		users: make(map[uint64]*User),
		mtx:   &sync.Mutex{},
	}
}

func (aois *AOIS) AddUser(u *User) {
	aois.mtx.Lock()
	defer aois.mtx.Unlock()
	aois.users[u.EntityID] = u
}

func (aois *AOIS) GetUser(entityID uint64) (*User, bool) {
	aois.mtx.Lock()
	defer aois.mtx.Unlock()
	u, ok := aois.users[entityID]
	return u, ok
}

func (aois *AOIS) AddAOI(flag byte, msg *msgdef.EnterAOI) {
	aois.mtx.Lock()
	defer aois.mtx.Unlock()

	log.Debug("aoi:", flag, " ", msg.EntityID)

	if flag == 1 {
		user, ok := aois.users[msg.EntityID]
		if !ok {
			user = NewUser()
			user.EntityID = msg.EntityID
			user.EntityType = msg.EntityType
			aois.users[msg.EntityID] = user
		}
		user.ParseProps(msg.PropNum, msg.Properties)
		user.SetPos(msg.Pos.Mul(MapRate))
		user.SetRota(msg.Rota)
	} else {
		if u, ok := aois.users[msg.EntityID]; ok {
			u.Destroy()
			delete(aois.users, msg.EntityID)
		}
	}
}

func (aois *AOIS) GetUserViews() []*User {
	aois.mtx.Lock()
	defer aois.mtx.Unlock()

	var lst []*User
	for _, v := range aois.users {
		lst = append(lst, v)
	}
	return lst
}

func (aois *AOIS) AOISLoop() {
	aois.mtx.Lock()
	defer aois.mtx.Unlock()

	for _, v := range aois.users {
		v.Loop()
	}
}
