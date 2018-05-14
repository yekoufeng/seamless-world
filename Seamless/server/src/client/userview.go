package main

import (
	"fmt"
	"math"
	"protoMsg"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type UserView struct {
	props    Props
	us       UserState
	um       *UserModel
	entityID uint64
	selected bool

	effects   []*protoMsg.EffectNotify
	effectMtx sync.Mutex
}

func (uv *UserView) CloneData(u *User) {
	u.CloneProps(&uv.props)
	u.CloneState(&uv.us)
	u.entityID = u.EntityID
}

func (uv *UserView) DrawModel(renderer *sdl.Renderer) {
	texture, rect := uv.um.model.GetFrame()
	if texture == nil {
		return
	}

	flip := sdl.FLIP_NONE
	angle := uv.us.GetAngle()
	if angle > math.Pi {
		flip = sdl.FLIP_HORIZONTAL
	}

	// log.Debug("angle:", angle, " with rota:", uv.us.rota)
	// seelog.Debug("draw at ", x, " ", z)
	showRect := uv.GetShowRect()
	renderer.CopyEx(texture, rect, showRect, 0, &sdl.Point{0, 0}, flip)

	//血条
	hpRect := &sdl.Rect{showRect.X, showRect.Y - 10, showRect.W, 10}

	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.FillRect(hpRect)

	rate := float32(uv.props.Hp) / float32(uv.props.MaxHp)
	curRect := &sdl.Rect{hpRect.X, hpRect.Y, int32(float32(hpRect.W) * rate), hpRect.H}

	renderer.SetDrawColor(102, 205, 0, 255)
	renderer.FillRect(curRect)

	renderer.SetDrawColor(69, 139, 0, 255)
	renderer.DrawRect(hpRect)

	//名字
	if uv.entityID == GetClient().EntityID {
		GetMainView().DrawText(hpRect.X, hpRect.Y-20, uv.props.Name, 0, 255, 0)
	} else {
		GetMainView().DrawText(hpRect.X, hpRect.Y-20, uv.props.Name, 255, 0, 0)
	}
}

func (uv *UserView) GetShowRect() *sdl.Rect {
	_, rect := uv.um.model.GetFrame()
	x, z := uv.us.GetUserViewPos()
	xOffset, zOffset := uv.um.model.GetOffset()
	x -= float32(xOffset)
	z -= float32(zOffset)
	return &sdl.Rect{int32(x), int32(z), rect.W, rect.H}
}

func (uv *UserView) GetRect() *sdl.Rect {
	x, y := uv.us.GetUserViewPos()
	xOffset, yOffset := uv.um.model.GetOffset()
	uv.um.rect.X = int32(x) - xOffset
	uv.um.rect.Y = int32(y) - yOffset
	// seelog.Debug("GetRect:", uv.um.rect)
	return &uv.um.rect
}

func (uv *UserView) OnClick(x, y int32) {
	// seelog.Debug("OnClick:", uv.entityID)
	if GetClient().EntityID == uv.entityID {
		return
	}
	GetClient().SetCurUser(uv)
	uv.selected = true
}

func (uv *UserView) OnMouseMove(x, y int32) {

}

func (uv *UserView) OnLostFocus() {
	GetClient().ClearIfIsMe(uv)
	uv.selected = false
}

func (uv *UserView) Draw(render *sdl.Renderer) {
	if uv.selected {
		render.SetDrawColor(188, 238, 104, 255)
		render.DrawRect(uv.GetShowRect())
	}
}

func (uv *UserView) AddEffect(effect *protoMsg.EffectNotify) {
	uv.effectMtx.Lock()
	defer uv.effectMtx.Unlock()
	uv.effects = append(uv.effects, effect)
}

func (uv *UserView) DrawEffect() {
	uv.effectMtx.Lock()
	defer uv.effectMtx.Unlock()
	for _, effect := range uv.effects {
		x, z := uv.us.GetUserViewPos()
		xOffset, zOffset := uv.um.model.GetOffset()
		x -= float32(xOffset)
		z -= float32(zOffset)
		text, r, g, b := uv.getEffectNotifyInfo(effect)
		t := NewUserEffectText(uv, text, r, g, b, sdl.Point{int32(x), int32(z)}, sdl.Point{int32(x), int32(z) - 100}, 2000)
		GetTransformMgr().AddTransform(t)
	}
	if len(uv.effects) > 0 {
		uv.effects = uv.effects[0:0]
	}
}

func (uv *UserView) getEffectNotifyInfo(effect *protoMsg.EffectNotify) (text string, r, g, b uint8) {
	switch effect.EffectType {
	case protoMsg.EffectType_Damage:
		text = fmt.Sprintf("-%v", effect.EffectParam)
		r = 255
	case protoMsg.EffectType_RecoverHp:
		text = fmt.Sprintf("+%v", effect.EffectParam)
		g = 255
	case protoMsg.EffectType_ReduceDefence:
		text = fmt.Sprintf("Defence: -%v", effect.EffectParam)
		r = 255
	default:

	}
	return
}
