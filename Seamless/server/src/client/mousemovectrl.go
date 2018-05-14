package main

import (
	"sync/atomic"
	"time"
	"zeus/linmath"
)

func (us *UserState) MoveTo(viewPosX, viewPosZ int32) {
	if atomic.LoadInt32(&us.moveType) == MoveTp_Keyboard {
		return
	}

	us.mtx.Lock()
	defer us.mtx.Unlock()
	us.clickPos = us.ViewPosToMapPos(linmath.Vector3{float32(viewPosX), 0, float32(viewPosZ)})
	// seelog.Debug("clickPos:", us.clickPos)
	atomic.StoreInt32(&us.moveType, MoveTp_Mouse)
	atomic.StoreInt64(&us.lastMoveTime, time.Now().UnixNano()/1e6)
	GetClient().User.um.model.StartAnimate()

	us.rota = us.clickPos.Sub(us.pos)
	us.AddMask((0x8 | 0x10 | 0x20))
	//增加显示
	GetClickEffect().StartAnimate()
}

func (us *UserState) CancelMouseMove() {
	if atomic.CompareAndSwapInt32(&us.moveType, MoveTp_Mouse, MoveTp_None) {
		GetClient().User.um.model.StopAnimate()
		//取消显示
		GetClickEffect().StopAnimate()
	}
}

func (us *UserState) MoveCtrlByMouse() {
	us.mtx.Lock()
	clickPos := us.clickPos
	us.mtx.Unlock()

	us.mtx.Lock()
	defer us.mtx.Unlock()

	delta := clickPos.Sub(us.pos)
	// seelog.Debug("delta1:", delta)
	l := delta.Len()
	if l <= 1 {
		atomic.StoreInt32(&us.moveType, MoveTp_None)
		GetClient().User.um.model.StopAnimate()
		GetClickEffect().StopAnimate()
		GetClient().SendStopMoveMsg()
		return
	}
	now := time.Now().UnixNano() / 1e6
	totalTime := l / float32(us.speed) * 1000
	lastMoveTime := atomic.LoadInt64(&us.lastMoveTime)
	// seelog.Debug("now:", now, " lastMoveTime:", lastMoveTime, " totalTime:", totalTime)
	delta = delta.Mul(float32(now-lastMoveTime) / totalTime)
	// seelog.Debug("delta2:", delta)
	us.pos = us.pos.Add(delta)
	us.VerifyPos()
	us.AddMask((0x1 | 0x2 | 0x4))
	atomic.StoreInt64(&us.lastMoveTime, now)
}
