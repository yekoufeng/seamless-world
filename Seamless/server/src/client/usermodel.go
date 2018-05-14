package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type UserModel struct {
	model IAnimate
	rect  sdl.Rect
}

func NewUserModel(render *sdl.Renderer) *UserModel {
	um := &UserModel{
		model: NewAnimateImages(render, "assets/user/%d.png", 8, 200, 30, 40),
	}

	_, surface := um.model.GetFrame()
	um.rect.W = surface.W
	um.rect.H = surface.H

	return um
}
