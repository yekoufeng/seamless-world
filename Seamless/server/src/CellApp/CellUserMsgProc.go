package main

import (
	"protoMsg"
	"zeus/iserver"
	"zeus/linmath"

	log "github.com/cihub/seelog"
)

// CellUserMsgProc CellUser的消息处理函数
type CellUserMsgProc struct {
	user *CellUser
}

// RPC_DoSkill
func (p *CellUserMsgProc) RPC_DoSkill(msg *protoMsg.AttackReq) {
	skill, ok := p.user.GetSkill(msg.SkillID)
	if !ok {
		log.Error("skill is not exist, skillid:", msg.SkillID)
		return
	}

	//TODO:各种条件判断

	if msg.TargetID != 0 {
		if target, ok := p.user.GetCell().GetEntity(msg.TargetID).(*CellUser); ok {
			skill.DoSkill(p.user, target, linmath.NewVector3(msg.X, 0, msg.Z))
		}
	} else {
		skill.DoSkill(p.user, nil, linmath.NewVector3(msg.X, 0, msg.Z))
	}
}

func (p *CellUserMsgProc) MsgProc_SkillEffect(msg *protoMsg.SkillEffect) {
	if p.user.IsGhost() {
		p.user.SendMsgToReal(msg)
		return
	}

	if msg.Recoverhp > 0 {
		p.user.AddHp(msg.Recoverhp)
	}
	if msg.Damage > 0 {
		p.user.MinusHp(msg.Damage)
	}

	if msg.ReduceDefence > 0 {
		p.user.MinusDefence(msg.ReduceDefence)
	}
	if msg.AddDefence > 0 {
		p.user.AddDefence(msg.AddDefence)
	}

	for _, effectid := range msg.Bufflist {
		if effect, ok := GetEffectCfgMgr().GetEffect(effectid); ok {
			effect.EffectToUser(p.user)
			p.user.Error("RPC_SkillEffect.effect:", effectid)
		}
	}
	//TODO:其它变化
}

func (p *CellUserMsgProc) MsgProc_DetectCell(msg *protoMsg.DetectCell) {
	// p.user.Debug("MsgProc_DetectCell")
	p.user.Post(iserver.ServerTypeClient, msg)
}
