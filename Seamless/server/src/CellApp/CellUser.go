package main

import (
	"common"
	"entitydef"
	"protoMsg"
	"zeus/dbservice"
	"zeus/linmath"

	log "github.com/cihub/seelog"
)

// CellUser
type CellUser struct {
	entitydef.PlayerDef
	Entity

	fightTmpData protoMsg.SkillEffect
	skillMap     map[uint32]*Skill

	timerMgr *common.TimerMgr
	buffMgr  *BuffMgr
}

// Init 初始化调用
func (user *CellUser) Init(initParam interface{}) {
	user.RegMsgProc(&CellUserMsgProc{user: user}) //注册消息处理对象

	user.SetWatcher(true)

	user.skillMap = make(map[uint32]*Skill)

	user.timerMgr = common.NewTimerMgr()

	user.buffMgr = NewBuffMgr(user)

	pos, ok := initParam.(*linmath.Vector3)
	if ok {
		user.SetPos(*pos)
	}

	user.testInit()
}

// testInit 测试数据
func (user *CellUser) testInit() {
	user.SetMaxHP(100)
	user.SetHP(100)
	user.SetAttack(20)
	user.SetDefence(5)
	// user.SetPos(linmath.Vector3{50, 0, 50})

	//理论上技能这一块是根据配置或者其它系统来初始化的
	//普攻带20%回血
	skill := &Skill{}
	effect, _ := GetEffectCfgMgr().GetEffect(101)
	skill.AddEffect(effect)
	effect, _ = GetEffectCfgMgr().GetEffect(102)
	skill.AddEffect(effect)
	user.skillMap[Skill_NormalAttack] = skill

	skill = &Skill{}
	effect, _ = GetEffectCfgMgr().GetEffect(103)
	skill.AddEffect(effect)
	user.skillMap[Skill_ReduceDefence] = skill

	if user.GetName() == "" {
		name, _ := dbservice.Account(user.GetDBID()).GetUsername()
		user.SetName(name)
	}
}

//Loop 定时执行
func (user *CellUser) Loop() {
	user.timerMgr.Tick()
}

func (user *CellUser) loadData() {
}

// Destroy 析构时调用
func (user *CellUser) Destroy() {
	user.GetEntities().UnregTimerByObj(user)
	if user.GetCell() != nil {
		user.GetCell().OnEntityDestory(&user.Entity)
	}
}

// OnEnterSpace 玩家进入地图
func (user *CellUser) OnEnterCell() {
	log.Debug("进入地图 ", user)

	//临时发送cellinfos给客户端
	user.SendCellInfos()
	//user.SendFullAOIs()
	user.SendFullProps()
}

// OnLeaveCell 玩家离开地图
func (user *CellUser) OnLeaveCell() {
	user.Info("离开地图", user)
}

func (user *CellUser) GetEntity() *Entity {
	return &user.Entity
}

func (user *CellUser) AddHp(hp uint32) uint32 {
	if user.GetHP() == 0 { //已经挂了
		return 0
	}
	recover := hp
	if recover+user.GetHP() > user.GetMaxHP() {
		recover = user.GetMaxHP() - user.GetHP()
	}
	user.SetHP(user.GetHP() + recover)
	return recover
}

func (user *CellUser) MinusHp(hp uint32) {
	if hp > user.GetHP() {
		hp = user.GetHP()
	}
	user.SetHP(user.GetHP() - hp)
}

func (user *CellUser) MinusDefence(defence uint32) {
	if defence > user.GetDefence() {
		defence = user.GetDefence()
	}
	user.SetDefence(user.GetDefence() - defence)
}

func (user *CellUser) AddDefence(defence uint32) {
	user.SetDefence(user.GetDefence() + defence)
}

func (user *CellUser) GetSkill(id uint32) (*Skill, bool) {
	v, ok := user.skillMap[id]
	return v, ok
}

func (user *CellUser) AddState(stateMask int) {
	state := user.GetState() | uint32(stateMask)
	user.SetState(state)
}

func (user *CellUser) RemoveState(stateMask int) {
	state := user.GetState() & (^uint32(stateMask))
	user.SetState(state)
}

func (user *CellUser) CheckState(stateMask int) bool {
	return (user.GetState() & uint32(stateMask)) > 0
}

// OnPosChange 玩家位置改变
func (user *CellUser) OnPosChange(curPos, origPos linmath.Vector3) {
	// if user.stateMgr == nil {
	// 	return
	// }

	// 精度1位，视为位置不变
	if int(curPos.X*10) == int(origPos.X*10) &&
		int(curPos.Y*10) == int(origPos.Y*10) &&
		int(curPos.Z*10) == int(origPos.Z*10) {
		return
	}

	// log.Debug("玩家位置改变 originPos:", origPos, " curPos:", curPos, "user id:", user.GetID())

}
