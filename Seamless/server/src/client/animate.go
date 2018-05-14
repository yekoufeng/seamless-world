package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type IAnimate interface {
	GetFrame() (*sdl.Texture, *sdl.Rect)

	StartAnimate()
	StopAnimate()

	GetOffset() (int32, int32)
	IsPlaying() bool
}

type AnimateBase struct {
	frames       uint32
	path         string
	interval     int64
	lastShowTime int64

	xOffset int32
	zOffset int32
}

func (a *AnimateBase) GetOffset() (int32, int32) {
	return a.xOffset, a.zOffset
}

func (a *AnimateBase) StopAnimate() {
	if a == nil {
		return
	}
	atomic.StoreInt64(&a.lastShowTime, 0)
}

func (a *AnimateBase) StartAnimate() {
	if a == nil || atomic.LoadInt64(&a.lastShowTime) > 0 {
		return
	}
	atomic.StoreInt64(&a.lastShowTime, time.Now().UnixNano()/1e6)
}

func (a *AnimateBase) IsPlaying() bool {
	return atomic.LoadInt64(&a.lastShowTime) > 0
}

type AnimateImages struct {
	AnimateBase

	textures []*ImagePair
}

func NewAnimateImages(render *sdl.Renderer, path string, frames uint32, interval int64, xOffset, zOffset int32) IAnimate {
	a := &AnimateImages{
		AnimateBase: AnimateBase{
			frames:   frames,
			path:     path,
			interval: interval,
			xOffset:  xOffset,
			zOffset:  zOffset,
		},
	}

	for i := 1; i <= int(a.frames); i++ {
		pair := GetAssetsMgr().GetImagePair(fmt.Sprintf(path, i))
		a.textures = append(a.textures, pair)
	}

	return a
}

func (a *AnimateImages) GetFrame() (*sdl.Texture, *sdl.Rect) {
	lastShowTime := atomic.LoadInt64(&a.lastShowTime)
	var frame *ImagePair
	if lastShowTime == 0 {
		frame = a.textures[0]
	} else {
		index := (time.Now().UnixNano()/1e6 - lastShowTime) / a.interval
		frame = a.textures[int(uint32(index)%a.frames)]
	}
	return frame.texture, &sdl.Rect{0, 0, frame.surface.W, frame.surface.H}
}

type AnimateOneImage struct {
	AnimateBase

	surface *sdl.Surface
	texture *sdl.Texture
}

func NewAnimateOneImage(render *sdl.Renderer, path string, frames uint32, interval int64, xOffset, zOffset int32) IAnimate {
	a := &AnimateOneImage{
		AnimateBase: AnimateBase{
			frames:   frames,
			path:     path,
			interval: interval,
			xOffset:  xOffset,
			zOffset:  zOffset,
		},
	}

	pair := GetAssetsMgr().GetImagePair(a.path)
	a.surface, a.texture = pair.surface, pair.texture

	return a
}

func (a *AnimateOneImage) GetFrame() (*sdl.Texture, *sdl.Rect) {
	lastShowTime := atomic.LoadInt64(&a.lastShowTime)
	if lastShowTime == 0 {
		return a.texture, &sdl.Rect{0, 0, a.surface.W / 8, a.surface.H}
	}
	index := (time.Now().UnixNano()/1e6 - lastShowTime) / a.interval
	return a.texture, &sdl.Rect{a.surface.W / 8 * int32(uint32(index)%a.frames), 0, a.surface.W / 8, a.surface.H}
}
