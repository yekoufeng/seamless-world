package main

import (
	"math"
	"sync/atomic"
	"time"
)

func (us *UserState) MoveCtrlByKeyboard() {
	us.mtx.Lock()
	defer us.mtx.Unlock()

	angle := us.GetAngle()

	lastMoveTime := atomic.LoadInt64(&us.lastMoveTime)
	// log.Debug("move with rota:", us.rota, " angle:", angle)
	now := time.Now().UnixNano() / 1e6
	us.pos.X += float32(us.speed*math.Sin(angle)) * float32(now-lastMoveTime) / 1000
	us.pos.Z += float32(us.speed*math.Cos(angle)) * float32(now-lastMoveTime) / 1000
	us.VerifyPos()
	us.AddMask((0x1 | 0x2 | 0x4))
	atomic.StoreInt64(&us.lastMoveTime, now)
}

func (us *UserState) KeyDown(key int) {
	us.CancelMouseMove()
	us.keyMask |= key
	us.updateRotaAndMoveState()
	// log.Debug("key down ", key, " keymask:", us.keyMask)
}

func (us *UserState) KeyUp(key int) {
	us.keyMask &= ^key
	us.updateRotaAndMoveState()
	// log.Debug("key up ", key, " keymask:", us.keyMask)
}

func (us *UserState) updateRotaAndMoveState() {
	mask := us.keyMask
	if (mask&KeyMask_A) > 0 && (mask&KeyMask_D) > 0 {
		mask &= ^(KeyMask_A | KeyMask_D)
	}
	if (mask&KeyMask_W) > 0 && (mask&KeyMask_S) > 0 {
		mask &= ^(KeyMask_W | KeyMask_S)
	}

	us.mtx.Lock()
	defer us.mtx.Unlock()

	if mask == 0 {
		atomic.StoreInt32(&us.moveType, MoveTp_None)
		GetClient().User.um.model.StopAnimate()
		GetClient().SendStopMoveMsg()
		return
	}
	if atomic.LoadInt32(&us.moveType) == MoveTp_None {
		atomic.StoreInt64(&us.lastMoveTime, time.Now().UnixNano()/1e6)
		GetClient().User.um.model.StartAnimate()
		atomic.StoreInt32(&us.moveType, MoveTp_Keyboard)
	}

	if mask == KeyMask_W {
		us.rota.X = 0
		us.rota.Z = -1
	} else if mask == (KeyMask_W | KeyMask_D) {
		us.rota.X = 1
		us.rota.Z = -1
	} else if mask == KeyMask_D {
		us.rota.X = 1
		us.rota.Z = 0
	} else if mask == (KeyMask_S | KeyMask_D) {
		us.rota.X = 1
		us.rota.Z = 1
	} else if mask == KeyMask_S {
		us.rota.X = 0
		us.rota.Z = 1
	} else if mask == (KeyMask_S | KeyMask_A) {
		us.rota.X = -1
		us.rota.Z = 1
	} else if mask == KeyMask_A {
		us.rota.X = -1
		us.rota.Z = 0
	} else if mask == (KeyMask_A | KeyMask_W) {
		us.rota.X = -1
		us.rota.Z = -1
	}
	us.AddMask((0x8 | 0x10 | 0x20))
	// log.Debug("rota:", us.rota)
}
