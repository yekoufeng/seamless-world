package main

import (
	"protoMsg"
	"time"
	"zeus/iserver"
	"zeus/linmath"
)

/*
计划实现3个效果
1. 攻击时根据造成的伤害百分比回血
2. 攻击时减少目标防御
3. 无敌
*/

const (
	Effect_Damage1           = 1 //根据攻防计算伤害
	Effect_RecoverHpByDamage = 2 //攻击时根据造成的伤害百分比回血
	Effect_ReduceDefence     = 3 //攻击时减少目标防御
	Effect_Unbeatable        = 4 //无敌
)

type IEffect interface {
	GetID() uint32
	GetType() uint32
	DoEffect(attacker, defencer *CellUser, pos linmath.Vector3, ulst map[uint64]*CellUser)
	EffectToUser(user *CellUser)
}

type Effect struct {
	ID   uint32
	Type uint32
}

func (effect *Effect) GetID() uint32 {
	return effect.ID
}

func (effect *Effect) GetType() uint32 {
	return effect.Type
}

func (effect *Effect) EffectToUser(user *CellUser) {

}

type EffectDamage1 struct {
	Effect
}

func (effect *EffectDamage1) DoEffect(attacker, defencer *CellUser, pos linmath.Vector3, ulst map[uint64]*CellUser) {
	if defencer == nil || defencer.GetHP() == 0 { // 已经死了就返回
		return
	}
	if defencer.CheckState(CellUserState_UnBeatable) {
		return
	}
	damage := attacker.GetAttack() - defencer.GetDefence()
	if defencer.GetHP() < damage {
		damage = defencer.GetHP()
	}
	defencer.SetHP(defencer.GetHP() - damage)
	defencer.fightTmpData.Damage += damage
	//广播攻击效果、伤害
	attacker.Debug("be attack lost hp:", damage)
	attacker.CastMsgToAllClient(&protoMsg.EffectNotify{
		EntityID:    defencer.GetID(),
		EffectType:  protoMsg.EffectType_Damage,
		EffectParam: int64(damage),
	})
	ulst[defencer.GetID()] = defencer
}

type EffectRecoverHpByDamage struct {
	Effect
	Rate float32 //回血百分比
}

func (effect *EffectRecoverHpByDamage) DoEffect(attacker, defencer *CellUser, pos linmath.Vector3, ulst map[uint64]*CellUser) {
	if attacker.GetHP() == 0 { //死亡就不能回了
		return
	}
	if defencer == nil {
		return
	}
	recoverhp := uint32(float32(defencer.fightTmpData.Damage) * effect.Rate)
	realAdd := attacker.AddHp(recoverhp)
	attacker.Debug("attack recover hp:", recoverhp)
	//广播攻击效果、伤害
	attacker.CastMsgToAllClient(&protoMsg.EffectNotify{
		EntityID:    attacker.GetID(),
		EffectType:  protoMsg.EffectType_RecoverHp,
		EffectParam: int64(realAdd),
	})
	ulst[attacker.GetID()] = attacker
}

type IEffectBuff interface {
	IEffect
	GetDuration() int64
}
type EffectBuff struct {
	Effect
	Duration int64 //持续时间 ms
}

func (effect *EffectBuff) GetDuration() int64 {
	return effect.Duration
}

type EffectReduceDefence struct {
	EffectBuff
	ReduceVal uint32
}

func (effect *EffectReduceDefence) DoEffect(attacker, defencer *CellUser, pos linmath.Vector3, ulst map[uint64]*CellUser) {
	cell := attacker.GetCell()
	if cell == nil {
		return
	}
	cell.TravsalAOI(attacker, func(o iserver.ICoordEntity) {
		if u, ok := cell.GetEntity(o.GetID()).(*CellUser); ok && u != attacker {
			deltaPos := o.GetPos().Sub(pos)
			if deltaPos.Len() < 2 {
				effect.EffectToUser(u)
				ulst[u.GetID()] = u
			}
		}
	})
}

func (effect *EffectReduceDefence) EffectToUser(user *CellUser) {
	if user == nil {
		return
	}
	if user.IsGhost() {
		user.fightTmpData.Bufflist = append(user.fightTmpData.Bufflist, effect.GetID())
		return
	}
	if user.GetHP() == 0 { //死亡就算了，绕过他
		return
	}

	if buff, ok := user.buffMgr.GetBuff(effect.ID); ok {
		user.buffMgr.RefreshTimer(buff)
		//TODO:广播攻击效果、伤害
		return
	}

	reduce := effect.ReduceVal
	if reduce > user.GetDefence() {
		reduce = user.GetDefence()
	}
	user.SetDefence(user.GetDefence() - reduce)

	user.CastMsgToAllClient(&protoMsg.EffectNotify{
		EntityID:    user.GetID(),
		EffectType:  protoMsg.EffectType_ReduceDefence,
		EffectParam: int64(reduce),
	})

	user.fightTmpData.ReduceDefence += reduce

	now := time.Now().UnixNano() / 1e6
	buff := NewBuff(now, effect, nil)
	user.buffMgr.AddBuff(buff)
	t := user.timerMgr.AddTimer(now+effect.Duration, func() {
		user.SetDefence(user.GetDefence() + reduce)
		user.buffMgr.RemoveBuff(buff)
		user.fightTmpData.AddDefence += reduce
		if user.IsGhost() {
			user.SendMsgToReal(&user.fightTmpData)
			user.Error("SendMsgToReal(&user.fightTmpData):", user.fightTmpData)
			user.fightTmpData.Reset()
		}
		user.CastMsgToAllClient(&protoMsg.EffectNotify{
			EntityID:    user.GetID(),
			EffectType:  protoMsg.EffectType_AddDefence,
			EffectParam: int64(reduce),
		})
	})
	buff.EndTimer = t
}

type EffectUnbeatable struct {
	EffectBuff
}

func (effect *EffectUnbeatable) DoEffect(attacker, defencer *CellUser, pos linmath.Vector3, ulst map[uint64]*CellUser) {
	effect.EffectToUser(attacker)
	ulst[attacker.GetID()] = attacker
}

func (effect *EffectUnbeatable) EffectToUser(user *CellUser) {
	// 这个效果就只给自己放吧，不考虑给队友
	if user == nil || user.GetHP() == 0 {
		return
	}
	if user.IsGhost() {
		user.fightTmpData.Bufflist = append(user.fightTmpData.Bufflist, effect.GetID())
		return
	}

	if buff, ok := user.buffMgr.GetBuff(effect.ID); ok {
		user.buffMgr.RefreshTimer(buff)
		//TODO:广播攻击效果、伤害
		return
	}

	now := time.Now().UnixNano() / 1e6
	//TODO:判断旧的无敌BUFF
	if user.buffMgr.UnBeatableBuff != nil {
		if user.buffMgr.UnBeatableBuff.IsBetter(now + effect.Duration) {
			//TODO:广播效果
			return
		}

		t := user.buffMgr.UnBeatableBuff.EndTimer
		if t != nil {
			user.timerMgr.Remove(t.GetIndex())
		}
	}

	user.SetState(CellUserState_UnBeatable)
	buff := NewBuff(now, effect, nil)
	user.buffMgr.AddBuff(buff)
	t := user.timerMgr.AddTimer(now+effect.Duration, func() {
		user.RemoveState(CellUserState_UnBeatable)
		if user.buffMgr.UnBeatableBuff == buff {
			user.buffMgr.UnBeatableBuff = nil
		}
		user.buffMgr.RemoveBuff(buff)
	})
	buff.EndTimer = t
	user.buffMgr.UnBeatableBuff = buff
}

// TODO:各个效果的实现
