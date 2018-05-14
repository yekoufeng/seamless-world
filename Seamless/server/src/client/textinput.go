package main

import (
	"strconv"
	"strings"
	"zeus/linmath"

	"github.com/veandco/go-sdl2/sdl"
)

type SimpleTextInput struct {
	rect sdl.Rect

	text  []byte
	focus bool
	limit int
}

func NewSimpleTextInput(x, y, w, h int32, limit int) *SimpleTextInput {
	return &SimpleTextInput{
		rect:  sdl.Rect{x, y, w, h},
		limit: limit,
	}
}

func (input *SimpleTextInput) GetRect() *sdl.Rect {
	return &input.rect
}
func (input *SimpleTextInput) OnClick(x, y int32) {
	input.focus = true
}
func (input *SimpleTextInput) OnLostFocus() {
	input.focus = false
}

func (input *SimpleTextInput) OnMouseMove(x, y int32) {

}

func (input *SimpleTextInput) Draw(render *sdl.Renderer) {
	// points := []sdl.Point{
	// 	sdl.Point{input.rect.X, input.rect.Y},
	// 	sdl.Point{input.rect.X + input.rect.W, input.rect.Y},
	// 	sdl.Point{input.rect.X + input.rect.W, input.rect.Y + input.rect.H},
	// 	sdl.Point{input.rect.X, input.rect.Y + input.rect.H},
	// 	sdl.Point{input.rect.X, input.rect.Y},
	// }
	render.SetDrawColor(255, 228, 225, 255)
	// render.DrawLines(points)
	render.DrawRect(&input.rect)

	if len(input.text) > 0 {
		GetMainView().DrawInputText(input.rect.X+5, input.rect.Y+3, string(input.text))
	}
}

// func (input *SimpleTextInput) UpdateClick(x, y int32) bool {
// 	p := &sdl.Point{x, y}
// 	input.focus = p.InRect(&input.rect)
// 	return input.focus
// }

func (input *SimpleTextInput) IsInputing() bool {
	return input.focus
}

func (input *SimpleTextInput) Append(key byte) {
	if len(input.text) >= input.limit {
		return
	}
	input.text = append(input.text, key)
}

func (input *SimpleTextInput) RemoveLast() {
	l := len(input.text)
	if l == 0 {
		return
	}
	input.text = input.text[0 : l-1]
}

func (input *SimpleTextInput) Do() {
	if len(input.text) == 0 {
		return
	}

	cmdstr := string(input.text)
	cmdargs := strings.Split(cmdstr, " ")
	l := len(cmdargs)
	if l == 0 {
		return
	}
	if cmdargs[0] == "moveto" && l == 3 {
		x, _ := strconv.Atoi(cmdargs[1])
		z, _ := strconv.Atoi(cmdargs[2])
		pos := linmath.Vector3{float32(x), 0, float32(z)}
		GetClient().SetPos(pos)
	}

	input.text = input.text[0:0]
}
