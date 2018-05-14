package main

import (
	"zeus/linmath"
)

const (
	Skill_NormalAttack  = 1001
	Skill_ReduceDefence = 1002
)

type Skill struct {
	// 按顺序触发效果
	effects []IEffect
}

func (skill *Skill) AddEffect(effect IEffect) {
	skill.effects = append(skill.effects, effect)
}

func (skill *Skill) DoSkill(attacker, defencer *CellUser, pos linmath.Vector3) {
	m := make(map[uint64]*CellUser)
	for _, effect := range skill.effects {
		effect.DoEffect(attacker, defencer, pos, m)
	}
	//如果是打在Ghost上的，那么把结果转给Real
	for _, u := range m {
		if u.IsGhost() {
			u.SendMsgToReal(&u.fightTmpData)
			u.Error("SendMsgToReal(&defencer.fightTmpData):", u.fightTmpData)
			u.fightTmpData.Reset()
		}
	}
}
