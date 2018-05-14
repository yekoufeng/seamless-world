package main

import (
	"math"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	_FlySpd = 1000
)

type SkillMgr struct {
	showSkill1Area bool
}

func (sm *SkillMgr) ShowNormalAttack() {
	uv := GetClient().GetCurUser()
	if uv == nil {
		return
	}
	rota := uv.us.pos.Sub(GetClient().UserView.us.pos)
	angle := GetVecAngle(&rota)
	rect := &sdl.Rect{0, 0, 20, 100}
	// xOffset, yOffset := GetClient().User.UserView.um.model.GetOffset()
	x, y := GetClient().User.UserView.us.GetUserViewPos()
	// angle := GetClient().User.UserView.us.GetAngle()
	angle2 := math.Pi*2 - angle
	// if angle2 < 0 {
	// 	angle2 += math.Pi * 2
	// }
	// seelog.Debug("NormalAttack angle:", angle, " rota:", GetClient().User.UserView.us.rota)
	t := NewImagePosTransform("assets/attack.png",
		sdl.Point{10, 0},
		angle2,
		rect,
		sdl.Point{int32(x), int32(y)},
		sdl.Point{int32(x + rota.X), int32(y + rota.Z)},
		int64((rota.Len()*1000))/int64(_FlySpd))
	GetTransformMgr().AddTransform(t)

	//发消息给服务器
	GetClient().NormalAttack(uv.entityID)
}

func (sm *SkillMgr) ShowSkill1() {
	sm.showSkill1Area = true
}
func (sm *SkillMgr) CancelShowSkill1() {
	sm.showSkill1Area = false
}

func (sm *SkillMgr) DrawSkillArea(render *sdl.Renderer) {
	if !sm.showSkill1Area {
		return
	}
	// ip := GetAssetsMgr().GetImagePair("assets/circle.png")
	p := GetMainView().GetMousePos()
	// rect := &sdl.Rect{p.X - 100, p.Y - 100, 200, 200}
	// render.Copy(ip.texture, nil, rect)

	gfx.CircleRGBA(render, p.X, p.Y, 100, 136, 206, 250, 255)
	gfx.FilledCircleRGBA(render, p.X, p.Y, 100, 136, 206, 250, 80)
}

func (sm *SkillMgr) TryDoSkill() bool {
	if sm.showSkill1Area {
		GetClient().DoSkill1()
		return true
	}
	return false
}
