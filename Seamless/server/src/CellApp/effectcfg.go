package main

type EffectCfgMgr struct {
	effectMap map[uint32]IEffect
}

var effectCfgMgr *EffectCfgMgr

func GetEffectCfgMgr() *EffectCfgMgr {
	if effectCfgMgr == nil {
		effectCfgMgr = &EffectCfgMgr{
			effectMap: make(map[uint32]IEffect),
		}
	}
	return effectCfgMgr
}

func (mgr *EffectCfgMgr) GetEffect(id uint32) (IEffect, bool) {
	effect, ok := mgr.effectMap[id]
	return effect, ok
}

func init() {
	//随便加些数据
	GetEffectCfgMgr().effectMap[101] = &EffectDamage1{
		Effect: Effect{
			ID:   101,
			Type: Effect_Damage1,
		},
	}
	GetEffectCfgMgr().effectMap[102] = &EffectRecoverHpByDamage{
		Effect: Effect{
			ID:   102,
			Type: Effect_RecoverHpByDamage,
		},
		Rate: 0.20,
	}
	GetEffectCfgMgr().effectMap[103] = &EffectReduceDefence{
		EffectBuff: EffectBuff{
			Effect: Effect{
				ID:   103,
				Type: Effect_ReduceDefence,
			},
			Duration: 10000,
		},
		ReduceVal: 3,
	}

	GetEffectCfgMgr().effectMap[104] = &EffectUnbeatable{
		EffectBuff: EffectBuff{
			Effect: Effect{
				ID:   104,
				Type: Effect_Unbeatable,
			},
			Duration: 5000,
		},
	}
}
