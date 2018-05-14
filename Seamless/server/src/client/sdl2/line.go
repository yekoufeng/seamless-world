package sdl2

import (
	"container/list"

	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type Point2D struct {
	X int32
	Y int32
}

type Line struct {
	Color      [4]uint8
	StartPoint *Point2D
	EndPoint   *Point2D
}

func GetNewLine(r, g, b, a uint8, startX, startY, endX, endY int32) *Line {
	var obj = &Line{}
	obj.Color[0] = r
	obj.Color[1] = g
	obj.Color[2] = b
	obj.Color[3] = a

	obj.StartPoint.X = startX
	obj.StartPoint.Y = startY
	obj.EndPoint.X = endX
	obj.EndPoint.Y = endY

	return obj
}

func (l *Line) Draw(renderer *sdl.Renderer) {
	renderer.SetDrawColor(l.Color[0], l.Color[1], l.Color[2], l.Color[3])
	rect := sdl.Rect{
		X:l.StartPoint.X,
		Y:l.StartPoint.Y,
		W:l.EndPoint.X - l.StartPoint.X,
		H:l.EndPoint.Y - l.EndPoint.Y,
		}
	renderer.DrawRect(&rect)
}

func (l *Line) Destroy() {
	// Todo： 不做任何事, 以后可做 Interface
}

// 颜色款式相同个的线
type LineMgr struct {
	lines    *list.List
	renderer *sdl.Renderer
}

var gLineMgr *LineMgr
var once sync.Once // 保证只执行一次
func GetNewLineMgr() *LineMgr {
	once.Do(func() {
		gLineMgr = &LineMgr{
			lines: list.New(),
		}
	})
	return gLineMgr
}

func (lm *LineMgr) Init(renderer *sdl.Renderer) {
	lm.renderer = renderer
}

func (lm *LineMgr) AddList(args *list.List) {
	lm.lines.PushBackList(args)
}

func (lm *LineMgr) AddOneMore(args ...*Line) {
	for _, line := range args {
		for e := lm.lines.Front(); e != nil; e = e.Next() {
			ptrline := e.Value.(*Line)
			if ptrline.StartPoint == line.StartPoint && ptrline.EndPoint == line.EndPoint {
				//fmt.Println("LineMgar::AddOne current line is equal, ignor...")
				return
			}
		}

		lm.lines.PushBack(line)
	}

}

func (lm *LineMgr) Draw() {
	for e := lm.lines.Front(); e != nil; e = e.Next() {
		obj := e.Value.(*Line)
		obj.Draw(lm.renderer)
	}
}

func (lm *LineMgr) Destroy() {
	var next *list.Element
	for e := lm.lines.Front(); e != nil; e = next {
		next = e.Next()
		lm.lines.Remove(e)
	}
}
