package sdl2

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Wnd struct {
	Wdw    *sdl.Window
	Render *sdl.Renderer
}

func GetNewWindow(title string, x, y, w, h int32, flag uint32, index int, flag1 uint32) *Wnd {
	window, err := sdl.CreateWindow(title, x, y, w, h, flag)
	if err != nil {
		panic(err)
	}

	obj := &Wnd{
		Wdw:    window,
		Render: nil,
	}
	obj.CreateRenderer(index, flag1)
	return obj
}

func (wnd *Wnd) Destroy() {
	if wnd.Render != nil {
		wnd.Render.Destroy()
	} else {
		fmt.Println("struct Wnd Render is nil")
	}

	if wnd.Wdw != nil {
		wnd.Wdw.Destroy()
	} else {
		fmt.Println("Struct Wnd window is nil")
	}
}

func (wnd *Wnd) CreateRenderer(index int, flags uint32) {
	var err error
	if wnd.Render, err = sdl.CreateRenderer(wnd.Wdw, index, flags); err != nil {
		panic(err)
	}
}
