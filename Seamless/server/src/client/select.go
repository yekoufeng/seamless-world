package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type SelectEntity struct {
	curUser *UserView
}

func (se *SelectEntity) GetCurUser() *UserView {
	return se.curUser
}

func (se *SelectEntity) SetCurUser(uv *UserView) {
	se.curUser = uv
}

func (se *SelectEntity) ClearIfIsMe(uv *UserView) {
	if se.curUser == uv {
		se.curUser = nil
	}
}

func (se *SelectEntity) DrawSelect(render *sdl.Renderer) {
	if se.curUser == nil {
		return
	}

	render.SetDrawColor(188, 238, 104, 255)
	render.DrawRect(&sdl.Rect{MainViewWidth/2 - 124, 48, 244, 104})
	GetMainView().DrawText(MainViewWidth/2-120, 50, fmt.Sprintf("EntityID:%v", se.curUser.entityID), 255, 0, 0)
	GetMainView().DrawText(MainViewWidth/2-120, 75, fmt.Sprintf("Hp:%v", se.curUser.props.Hp), 255, 0, 0)
	GetMainView().DrawText(MainViewWidth/2-120, 100, fmt.Sprintf("Attack:%v", se.curUser.props.Attack), 255, 0, 0)
	GetMainView().DrawText(MainViewWidth/2-120, 125, fmt.Sprintf("Defence:%v", se.curUser.props.Defence), 255, 0, 0)
}
