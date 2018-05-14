package sdl2

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Text struct {
	Font *ttf.Font
}

func init() {
	if err := ttf.Init(); err != nil {
		panic(err)
	}
}

func GetNewText(path string, size int) *Text {
	/*
		path: ttf的路径
		size: 字体大小
	*/

	obj := &Text{}

	var err error
	if obj.Font, err = ttf.OpenFont(path, size); err != nil {
		panic(err)
	}

	return obj
}

func (tx *Text) Destroy() {
	if tx.Font == nil {
		fmt.Println("Struct Text Font is nil")
	}
	tx.Font.Close()
}

func (tx *Text) RenderUTF8Blended(input string, color sdl.Color, render *sdl.Renderer, wrapLength int, x, y int32) {
	var err error
	var Surface *sdl.Surface
	if Surface, err = tx.Font.RenderUTF8Blended(input, color); err != nil {
		panic(err)
	}
	defer Surface.Free()

	box := sdl.Rect{x, y, Surface.W, Surface.H}
	var texture *sdl.Texture
	if texture, err = render.CreateTextureFromSurface(Surface); err != nil {
		panic(err)
	}
	defer texture.Destroy()

	render.Copy(texture, nil, &box)
}
