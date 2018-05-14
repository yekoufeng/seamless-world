package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Tips struct {
	iconPath string
	rect     sdl.Rect
	showTips bool
	text     []string
	surface  *sdl.Surface
	texture  *sdl.Texture
}

func NewTips(render *sdl.Renderer, iconPath string, x, y, w, h int32, text []string) *Tips {
	t := &Tips{
		iconPath: iconPath,
		rect:     sdl.Rect{x, y, w, h},
		text:     text,
	}
	pair := GetAssetsMgr().GetImagePair(iconPath)
	t.surface, t.texture = pair.surface, pair.texture
	return t
}

func (t *Tips) GetRect() *sdl.Rect {
	return &t.rect
}

func (t *Tips) OnClick(x, y int32) {

}

func (t *Tips) OnLostFocus() {

}

func (t *Tips) OnMouseMove(x, y int32) {
	p := &sdl.Point{x, y}
	t.showTips = p.InRect(&t.rect)
}

func (t *Tips) Draw(render *sdl.Renderer) {
	render.Copy(t.texture, nil, &t.rect)

	if t.showTips {
		for i, s := range t.text {
			GetMainView().DrawText(MainViewWidth/2-100, 200+int32(i)*20, s, 0, 245, 255)
		}
	}
}
