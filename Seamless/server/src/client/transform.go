package main

import (
	"math"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/gfx"

	"github.com/veandco/go-sdl2/sdl"
)

type ITransform interface {
	Start()
	Update(now int64)
	Draw(render *sdl.Renderer)
	IsRuning() bool
}

type TransformBase struct {
	duration  int64 //多少ms完成
	startTime int64 //开始时间
}

func (tb *TransformBase) Start() {
	tb.startTime = time.Now().UnixNano() / 1e6
}

func (tb *TransformBase) IsRuning() bool {
	return time.Now().UnixNano()/1e6 < tb.duration+tb.startTime
}

type PosTransform struct {
	TransformBase
	start sdl.Point
	end   sdl.Point
	pos   sdl.Point
}

func (pt *PosTransform) Update(now int64) {
	r := float32(now-pt.startTime) / float32(pt.duration)
	pt.pos.X = int32(float32(pt.start.X) + float32(pt.end.X-pt.start.X)*r)
	pt.pos.Y = int32(float32(pt.start.Y) + float32(pt.end.Y-pt.start.Y)*r)
}

type ImagePosTransform struct {
	PosTransform
	ip     *ImagePair
	rect   sdl.Rect
	center sdl.Point
	angle  float64
}

func NewImagePosTransform(path string, center sdl.Point, angle float64, rect *sdl.Rect, start, end sdl.Point, duration int64) ITransform {
	// seelog.Debug("NewImagePosTransform")
	ip := GetAssetsMgr().GetImagePair(path)
	if rect == nil {
		rect = &sdl.Rect{0, 0, ip.surface.W, ip.surface.H}
	}
	t := &ImagePosTransform{
		PosTransform: PosTransform{
			TransformBase: TransformBase{
				duration: duration,
			},
			start: start,
			end:   end,
		},
		ip:     ip,
		rect:   *rect,
		center: center,
		angle:  angle * 360 / (math.Pi * 2),
	}
	return t
}

func (ipt *ImagePosTransform) Draw(render *sdl.Renderer) {
	// seelog.Debug("ImagePosTransform pos;", ipt.pos, " center:", ipt.center, " angle:", ipt.angle)
	render.CopyEx(ipt.ip.texture,
		nil,
		&sdl.Rect{ipt.pos.X, ipt.pos.Y, ipt.rect.W, ipt.rect.H},
		ipt.angle,
		&ipt.center,
		sdl.FLIP_NONE)
}

type TransformMgr struct {
	lst   []ITransform
	mutex sync.Mutex
}

var _TransformMgr *TransformMgr

func GetTransformMgr() *TransformMgr {
	if _TransformMgr == nil {
		_TransformMgr = &TransformMgr{}
	}
	return _TransformMgr
}

func (tm *TransformMgr) AddTransform(trans ITransform) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.lst = append(tm.lst, trans)
	trans.Start()
}

func (tm *TransformMgr) Update() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	now := time.Now().UnixNano() / 1e6
	for i := 0; i < len(tm.lst); {
		t := tm.lst[i]
		if !t.IsRuning() {
			for j := i; j < len(tm.lst)-1; j++ {
				tm.lst[j] = tm.lst[j+1]
			}
			tm.lst = tm.lst[0 : len(tm.lst)-1]
		} else {
			t.Update(now)
			i++
		}
	}
}

func (tm *TransformMgr) Draw(render *sdl.Renderer) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	for _, t := range tm.lst {
		t.Draw(render)
	}
}

type TextPosTransform struct {
	PosTransform
	text string
	r    uint8
	g    uint8
	b    uint8
}

func NewTextPosTransform(text string,
	r, g, b uint8,
	start, end sdl.Point,
	duration int64) *TextPosTransform {
	t := &TextPosTransform{
		PosTransform: PosTransform{
			TransformBase: TransformBase{
				duration: duration,
			},
			start: start,
			end:   end,
		},
		text: text,
		r:    r,
		g:    g,
		b:    b,
	}
	return t
}

func (tpt *TextPosTransform) Draw(render *sdl.Renderer) {
	gfx.StringRGBA(render, tpt.pos.X, tpt.pos.Y, tpt.text, tpt.r, tpt.g, tpt.b, 255)
}

type UserEffectText struct {
	TextPosTransform
	uv *UserView
}

func NewUserEffectText(uv *UserView,
	text string,
	r, g, b uint8,
	start, end sdl.Point,
	duration int64) *UserEffectText {
	t := &UserEffectText{
		TextPosTransform: TextPosTransform{
			PosTransform: PosTransform{
				TransformBase: TransformBase{
					duration: duration,
				},
				start: start,
				end:   end,
			},
			text: text,
			r:    r,
			g:    g,
			b:    b,
		},
		uv: uv,
	}
	return t
}

func (uet *UserEffectText) Draw(render *sdl.Renderer) {
	x, y := uet.uv.us.GetUserViewPos()
	xOffset, zOffset := uet.uv.um.model.GetOffset()
	x -= float32(xOffset)
	y -= float32(zOffset)
	xStart := int32(x) + uet.pos.X - uet.start.X
	yStart := int32(y) + uet.pos.Y - uet.start.Y
	gfx.StringRGBA(render, xStart, yStart, uet.text, uet.r, uet.g, uet.b, 255)
}
