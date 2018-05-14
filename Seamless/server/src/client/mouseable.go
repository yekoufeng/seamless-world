package main

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type IMouseable interface {
	GetRect() *sdl.Rect
	OnClick(x, y int32)
	OnMouseMove(x, y int32)
	OnLostFocus()
	Draw(*sdl.Renderer)
}

type MouseableMgr struct {
	lst   []IMouseable
	focus IMouseable
	mutex sync.Mutex
}

var _MouseableMgr *MouseableMgr

func GetMouseableMgr() *MouseableMgr {
	if _MouseableMgr == nil {
		_MouseableMgr = &MouseableMgr{}
	}
	return _MouseableMgr
}

func (mm *MouseableMgr) AddMouseable(i IMouseable) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.lst = append(mm.lst, i)
}

func (mm *MouseableMgr) RemoveMouseable(im IMouseable) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	for i := 0; i < len(mm.lst); {
		if mm.lst[i] == im {
			for j := i; j < len(mm.lst)-1; j++ {
				mm.lst[j] = mm.lst[j+1]
			}
			mm.lst = mm.lst[0 : len(mm.lst)-1]
			return
		} else {
			i++
		}
	}
}

func (mm *MouseableMgr) UpdateClick(x, y int32) bool {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	p := &sdl.Point{x, y}
	for _, m := range mm.lst {
		rect := m.GetRect()
		if p.InRect(rect) {
			if mm.focus != nil && m != mm.focus {
				mm.focus.OnLostFocus()
			}
			m.OnClick(x, y)
			mm.focus = m
			return true
		}
	}

	if mm.focus != nil {
		mm.focus.OnLostFocus()
	}
	return false
}

func (mm *MouseableMgr) MouseMove(x, y int32) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	for _, m := range mm.lst {
		m.OnMouseMove(x, y)
	}
}

func (mm *MouseableMgr) Draw(render *sdl.Renderer) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	for _, m := range mm.lst {
		m.Draw(render)
	}
}
