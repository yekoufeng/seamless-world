package main

import (
	"common"
	"time"
)

type IBuff interface {
	GetBeginTime() int64
	GetEffect() IEffectBuff
	GetEndTimer() *common.Timer
}

type Buff struct {
	BeginTime int64
	effect    IEffectBuff

	EndTimer *common.Timer
}

func (buff *Buff) GetBeginTime() int64 {
	return buff.BeginTime
}
func (buff *Buff) GetEffect() IEffectBuff {
	return buff.effect
}
func (buff *Buff) GetEndTimer() *common.Timer {
	return buff.EndTimer
}

func NewBuff(beginTime int64, effect IEffectBuff, t *common.Timer) *Buff {
	return &Buff{
		BeginTime: beginTime,
		effect:    effect,
		EndTimer:  t,
	}
}

func (buff *Buff) IsBetter(endTime int64) bool {
	return buff.BeginTime+buff.effect.GetDuration() >= endTime
}

type BuffMgr struct {
	user           *CellUser
	UnBeatableBuff *Buff

	BuffMap map[uint32]*Buff
}

func NewBuffMgr(user *CellUser) *BuffMgr {
	return &BuffMgr{
		user:    user,
		BuffMap: make(map[uint32]*Buff),
	}
}

func (bm *BuffMgr) AddBuff(buff *Buff) {
	bm.BuffMap[buff.effect.GetID()] = buff
}

func (bm *BuffMgr) RemoveBuff(buff *Buff) {
	delete(bm.BuffMap, buff.effect.GetID())
}

func (bm *BuffMgr) GetBuff(id uint32) (*Buff, bool) {
	buff, ok := bm.BuffMap[id]
	return buff, ok
}

func (bm *BuffMgr) RefreshTimer(buff *Buff) {
	if buff.EndTimer == nil {
		return
	}

	buff.EndTimer.SetTrigTime(time.Now().UnixNano()/1e6 + buff.effect.GetDuration())
	bm.user.timerMgr.Fix(buff.EndTimer.GetIndex())
}
