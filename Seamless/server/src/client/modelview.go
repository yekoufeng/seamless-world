package main

import "github.com/veandco/go-sdl2/sdl"

// log "github.com/cihub/seelog"

func (mv *MainView) DrawModel() {

	if GetClickEffect().IsPlaying() {
		texture, rect := GetClickEffect().GetFrame()
		pos := GetClient().UserView.us.clickPos
		mapRect := GetClient().UserView.us.GetViewMapRect()
		x := pos.X - float32(mapRect.X)
		z := pos.Z - float32(mapRect.Y)

		xOffset, zOffset := GetClickEffect().GetOffset()
		x -= float32(xOffset)
		z -= float32(zOffset)
		// log.Debug("angle:", angle, " with rota:", uv.us.rota)
		GetMainView().window.Render.Copy(texture, rect, &sdl.Rect{int32(x), int32(z), rect.W, rect.H})
	}

	lst := GetClient().AOIS.GetUserViews()
	if lst == nil {
		return
	}

	for _, v := range lst {
		v.CloneData(v)
		v.DrawModel(GetMainView().window.Render)
		v.DrawEffect()
	}
}
