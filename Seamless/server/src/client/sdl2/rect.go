package sdl2

import (
	"container/list"

	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type RectMgr struct {
	RectList *list.List
	Renderer *sdl.Renderer
}

var gRectMgr *RectMgr
var onceRect sync.Once // 保证只执行一次

func GetNewRectMgr() *RectMgr {
	onceRect.Do(func() {
		gRectMgr = &RectMgr{
			RectList: list.New(),
		}
	})
	return gRectMgr
}

func (rm *RectMgr) Init(renderer *sdl.Renderer) {
	rm.Renderer = renderer
}

func (rm *RectMgr) Add(args ...*sdl.Rect) {
	for _, rect := range args {
		rm.RectList.PushBack(rect)
	}
}

func (rm *RectMgr) Destroy() {
	var next *list.Element
	for e := rm.RectList.Front(); e != nil; e = next {
		next = e.Next()
		rm.RectList.Remove(e)
	}
}

func (rm *RectMgr) Draw() {
	rm.Renderer.SetDrawColor(255, 0, 0, 255)
	for e := rm.RectList.Front(); e != nil; e = e.Next() {
		rect := e.Value.(*sdl.Rect)
		rm.Renderer.DrawRect(rect)
	}
}

func (rm *RectMgr) DrawRectNow(r, g, b, a uint8, args ...*sdl.Rect) {
	rm.Renderer.SetDrawColor(r, g, b, a)
	for _, rect := range args {
		rm.Renderer.FillRect(rect)
		rm.Renderer.DrawRect(rect)
	}
}
