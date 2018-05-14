package main

import (
	"client/sdl2"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/veandco/go-sdl2/gfx"

	log "github.com/cihub/seelog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	MainViewWidth      = 1200
	MainViewHeight     = 900
	MainViewWrapLength = 200 // 限制字符屏幕显示的长度

	MapWidthWithPixel  = 5000
	MapHeightWithPixel = 5000

	MapWidth  = 100
	MapHeight = 100
	//jcmiao:坐标是M作单位，目前直接当像素用了，需要修改下

	MapRate = MapWidthWithPixel / MapWidth
)

const (
	// 地图相关信息
	MAPPATH = "assets/map"
	MAPNUM  = 4
)

type MainView struct {
	window *sdl2.Wnd
	text   *sdl2.Text
	mapmgr *sdl2.MapMgr
	stop   int32

	PropsView

	textInput *SimpleTextInput
	inputText *sdl2.Text

	enterCell    int32
	altDown      bool
	ctrlDown     bool
	keyZDown     bool
	mouseMovePos sdl.Point

	tabDown bool
}

func init() {

}

var _mainView *MainView

func GetMainView() *MainView {
	if _mainView == nil {
		_mainView = &MainView{
			stop: 0,
		}
		_mainView.PropsView = PropsView{
			mv: _mainView,
		}
		_mainView.Construct()
	}

	return _mainView
}

func (mv *MainView) Construct() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	// 创建窗口
	mv.window = sdl2.GetNewWindow("Seamless Client", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		MainViewWidth, MainViewHeight, sdl.WINDOW_SHOWN, -1, 0)
}

func (mv *MainView) Init() {
	// 创建字体
	mv.text = sdl2.GetNewText("consola.ttf", 20)
	mv.inputText = sdl2.GetNewText("consola.ttf", 30)

	// 创建背景
	mv.mapmgr = sdl2.GetNewMapMgr(MainViewWidth, MainViewHeight, MapWidthWithPixel, MapHeightWithPixel, mv.window.Render)
	mv.mapmgr.Load(MAPPATH, MAPNUM)

	//创建CellInfo矩形
	sdl2.GetNewRectMgr().Init(mv.window.Render)

	mv.textInput = NewSimpleTextInput(MainViewWidth/2-300, MainViewHeight-50, 600, 40, 35)
	GetMouseableMgr().AddMouseable(mv.textInput)

	GetMouseableMgr().AddMouseable(
		NewTips(mv.window.Render, "assets/question.jpg", MainViewWidth/2+310, MainViewHeight-50, 40, 40, GmTips))
}

func (mv *MainView) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Error(string(debug.Stack()))
		}
	}()

	for atomic.LoadInt32(&mv.stop) == 0 {
		select {
		case <-closeChan:
			atomic.CompareAndSwapInt32(&mv.stop, 0, 1)
		default:
		}
		if atomic.LoadInt32(&mv.enterCell) == 0 {
			mv.ShowEnterCelling()
		} else {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					log.Info("main view exit")
					atomic.CompareAndSwapInt32(&mv.stop, 0, 1)
					break
				case *sdl.KeyboardEvent:
					// 去除重复按键
					if event.(*sdl.KeyboardEvent).Repeat != 0 {
						continue
					}
					mv.KeyboardEvent(event.(*sdl.KeyboardEvent))
				case *sdl.MouseButtonEvent:
					mv.MouseButonEvent(event.(*sdl.MouseButtonEvent))
				case *sdl.MouseMotionEvent:
					mv.MouseMotionEvent(event.(*sdl.MouseMotionEvent))
				case *sdl.TextInputEvent:
					mv.TextInputEvent(event.(*sdl.TextInputEvent))
				default:
				}
			}
			//设置背景色
			GetClient().User.CloneData(GetClient().User)
			mv.window.Render.SetDrawColor(0, 0, 0, 255)
			mv.window.Render.Clear()
			mv.DrawBackground()
			mv.DrawCellRage()
			mv.DrawCellSvrID()
			mv.DrawCellInfo()
			mv.DrawModel()

			GetTransformMgr().Update()
			GetTransformMgr().Draw(mv.window.Render)

			mv.DrawProps()
			mv.DrawMousePos()
			mv.DrawCellSvrInfo()
			GetClient().DrawSelect(mv.window.Render)
			GetMouseableMgr().Draw(mv.window.Render)

			GetClient().DrawSkillArea(mv.window.Render)
			mv.DrawAllCellMiniMap()
		}

		mv.window.Render.Present()
		sdl.Delay(16)
	}

	mv.text.Destroy()
	mv.mapmgr.Destroy()
	mv.window.Destroy()
	ttf.Quit()
	sdl.Quit()

	//通知主goroutine退出
	ExitClient()
}

func (mv *MainView) SkillEvent(evt *sdl.KeyboardEvent) bool {
	if evt.Repeat != 0 {
		return false
	}
	if evt.Type == sdl.KEYDOWN {
		if evt.Keysym.Sym == sdl.K_1 {
			GetClient().ShowNormalAttack()
			return true
		} else if evt.Keysym.Sym == sdl.K_2 {
			GetClient().ShowSkill1()
			return true
		} else if evt.Keysym.Sym == sdl.K_ESCAPE {
			GetClient().CancelShowSkill1()
			return false
		}
	} else if evt.Type == sdl.KEYUP {

	}
	return false
}

func (mv *MainView) KeyboardEvent(evt *sdl.KeyboardEvent) {
	if evt.Keysym.Sym == sdl.K_LALT {
		mv.altDown = evt.Type == sdl.KEYDOWN
		return
	}
	if evt.Keysym.Sym == sdl.K_LCTRL {
		mv.ctrlDown = evt.Type == sdl.KEYDOWN
		return
	}
	if evt.Keysym.Sym == sdl.K_z {
		mv.keyZDown = evt.Type == sdl.KEYDOWN
		return
	}
	if evt.Keysym.Sym == sdl.K_TAB {
		mv.tabDown = evt.Type == sdl.KEYDOWN
		return
	}
	if mv.textInput.IsInputing() {
		if evt.Type == sdl.KEYDOWN {
			switch evt.Keysym.Sym {
			case sdl.K_BACKSPACE:
				mv.textInput.RemoveLast()
			case sdl.K_RETURN, sdl.K_KP_ENTER:
				mv.textInput.Do()
			default:
			}
		}
		return
	}

	if mv.SkillEvent(evt) {
		return
	}

	key := 0
	switch evt.Keysym.Sym {
	case sdl.K_w:
		key = KeyMask_W
		break
	case sdl.K_a:
		key = KeyMask_A
		break
	case sdl.K_s:
		key = KeyMask_S
		break
	case sdl.K_d:
		key = KeyMask_D
		break
	default:
		return
	}

	if evt.Type == sdl.KEYDOWN {
		GetClient().KeyDown(key)
	} else if evt.Type == sdl.KEYUP {
		GetClient().KeyUp(key)
	}
}

func (mv *MainView) MouseButonEvent(evt *sdl.MouseButtonEvent) {
	defer func() {
		GetClient().CancelShowSkill1()
	}()
	if evt.State == sdl.PRESSED {
		if evt.Button == sdl.BUTTON_LEFT {
			if GetClient().TryDoSkill() {
				return
			}
			if GetMouseableMgr().UpdateClick(evt.X, evt.Y) {
				return
			}
			// log.Debug("MouseButonEvent:BUTTON_LEFT")
			GetClient().MoveTo(evt.X, evt.Y)
		}
	} else {

	}
}

func (mv *MainView) MouseMotionEvent(evt *sdl.MouseMotionEvent) {
	GetMouseableMgr().MouseMove(evt.X, evt.Y)
	mv.mouseMovePos.X, mv.mouseMovePos.Y = evt.X, evt.Y
}

func (mv *MainView) DrawBackground() {
	rect := GetClient().UserView.us.GetViewMapRect()
	x := rect.X
	y := rect.Y
	var pos = [2]int32{int32(x), int32(y)}
	mv.mapmgr.Draw(pos)
}

func (mv *MainView) DrawText(x, y int32, text string, r, g, b uint8) {
	mv.text.RenderUTF8Blended(text, sdl.Color{r, g, b, 255}, mv.window.Render,
		MainViewWrapLength, x, y)
}

func (mv *MainView) DrawCellRage() {
	us := &GetClient().UserView.us
	rect := us.GetViewMapRect()
	var minX, minY int32
	minX = int32(rect.X)
	minY = int32(rect.Y)

	sdl2.GetNewRectMgr().Destroy()
	GetCellInfoMgr().GetDrawRect(minX, minY, MainViewWidth, MainViewHeight)
	sdl2.GetNewRectMgr().Draw()
}

func (mv *MainView) DrawCellSvrID() {
	us := &GetClient().UserView.us
	rect := us.GetViewMapRect()
	var minX, minY int32
	minX = int32(rect.X)
	minY = int32(rect.Y)

	GetCellInfoMgr().DrawCellID(minX, minY, MainViewWidth, MainViewHeight, mv)
}

func (mv *MainView) DrawInputText(x, y int32, text string) {
	mv.inputText.RenderUTF8Blended(text, sdl.Color{255, 0, 0, 255}, mv.window.Render,
		MainViewWrapLength, x, y)
}

func (mv *MainView) TextInputEvent(evt *sdl.TextInputEvent) {
	if mv.textInput.IsInputing() {
		mv.textInput.Append(evt.Text[0])
	}
}

func (mv *MainView) EnterCellOk() {
	atomic.StoreInt32(&mv.enterCell, 1)
}

func (mv *MainView) ShowEnterCelling() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	}
	mv.window.Render.SetDrawColor(0, 0, 0, 255)
	mv.window.Render.Clear()
	n := time.Now().Unix()%3 + 1
	str := []byte("enter cell ")
	for i := 0; i < int(n); i++ {
		_ = i
		str = append(str, '.')
	}
	mv.DrawText(MainViewWidth/2-80, MainViewHeight/2-20, string(str), 0, 245, 255)
	GetClient().TryEnterCell()
}

func (mv *MainView) DrawMousePos() {
	if !mv.altDown {
		return
	}
	x, y := mv.mouseMovePos.X, mv.mouseMovePos.Y
	rect := GetClient().User.UserView.us.GetViewMapRect()
	str := fmt.Sprintf("view(%v,%v) map(%v, %v)", x, y, rect.X+x, rect.Y+y)
	// x += 15
	if x > MainViewWidth-320 {
		x = MainViewWidth - 320
	}
	y -= 20
	if y > MainViewHeight-40 {
		y = MainViewHeight - 40
	}
	mv.DrawText(x, y, str, 0, 245, 255)
}

func (mv *MainView) DrawCellInfo() {
	if !mv.ctrlDown {
		return
	}
	if mv.tabDown {
		return
	}
	x, y := mv.mouseMovePos.X, mv.mouseMovePos.Y

	rect := GetClient().UserView.us.GetViewMapRect()
	// 显示当前鼠标点击的位置 获取相应Cell的信息
	cellinfo := GetCellInfoMgr().GetCellInfoByPos(x+rect.X, y+rect.Y)
	if rect := GetCellInfoMgr().GetRectByID(cellinfo.CellID); rect != nil {
		offsetPos := GetClient().UserView.us.GetViewMapRect()
		rect.X -= offsetPos.X
		rect.Y -= offsetPos.Y
		sdl2.GetNewRectMgr().DrawRectNow(255, 153, 0, 255, rect)
	}

	if cellinfo == nil {
		mv.DrawOne(x, y, "当前鼠标指向不在Cell范围内")
		return
	}
	mv.DrawOne(x, y-GridHeight, fmt.Sprintf("%10v-%10v", mv.mouseMovePos.X, mv.mouseMovePos.Y))
	mv.DrawOne(x, y, fmt.Sprintf("%10v-%10v", "SvrID: ", cellinfo.SrvID))
	mv.DrawOne(x, y+GridHeight, fmt.Sprintf("%10v-%10v", "SpaceID: ", cellinfo.SpaceID))
	mv.DrawOne(x, y+2*GridHeight, fmt.Sprintf("%10v-%10v", "CellID", cellinfo.CellID))
	mv.DrawOne(x, y+3*GridHeight, fmt.Sprintf("%10v-%10v", "CurMinX", cellinfo.MinX))
	mv.DrawOne(x, y+4*GridHeight, fmt.Sprintf("%10v-%10v", "CurMinY", cellinfo.MinY))
	mv.DrawOne(x, y+5*GridHeight, fmt.Sprintf("%10v-%10v", "CurMaxX", cellinfo.MaxX))
	mv.DrawOne(x, y+6*GridHeight, fmt.Sprintf("%10v-%10v", "CurMaxY", cellinfo.MaxY))

	mv.DrawOne(x, y+7*GridHeight, fmt.Sprintf("%10v-%10v", "SvrMinX", int32(cellinfo.MinX/MapRate)))
	mv.DrawOne(x, y+8*GridHeight, fmt.Sprintf("%10v-%10v", "SvrMinY", int32(cellinfo.MinY/MapRate)))
	mv.DrawOne(x, y+9*GridHeight, fmt.Sprintf("%10v-%10v", "SvrMaxX", int32(cellinfo.MaxX/MapRate)))
	mv.DrawOne(x, y+10*GridHeight, fmt.Sprintf("%10v-%10v", "SvrMaxY", int32(cellinfo.MaxY/MapRate)))

	//mv.DrawOne(0, 0, fmt.Sprintf("%10v%-10v", "SvrID: ", cellinfo.SrvID))
	//mv.DrawText(PVXBegin+(x+1)*GridWidth, PVYBegin+(y+1)*GridHeight,
	//	fmt.Sprintf("%10v%-10v", "SpaceID: ", cellinfo.SpaceID), 0, 245, 255)
}

func (mv *MainView) DrawCellSvrInfo() {
	if !mv.keyZDown {
		return
	}

	x, y := mv.mouseMovePos.X, mv.mouseMovePos.Y

	// 显示当前鼠标点击的位置 获取相应Cell的信息
	cellinfo := GetCellInfoMgr().GetAllCellInfo()

	index := 0
	for _, cell := range cellinfo {
		mv.DrawOne(x, y+int32(index)*GridHeight, fmt.Sprintf("SvrID:%5d CellID:%10d SpaceID:%10d MinX:%.3f MaxX:%.3f MinY:%.3f MaxY:%.3f", cell.SrvID, cell.CellID, cell.SpaceID, cell.MinX/MapRate, cell.MaxX/MapRate, cell.MinY/MapRate, cell.MaxY/MapRate))
		index += 1
	}
}

func (mv *MainView) DrawOne(x, y int32, text string) {
	mv.DrawText(x, PVYBegin+y, text, 0, 245, 255)
}

func (mv *MainView) GetMousePos() *sdl.Point {
	return &mv.mouseMovePos
}

func (mv *MainView) DrawAllCellMiniMap() bool {
	if !mv.tabDown {
		return false
	}

	// gfx.BoxRGBA(mv.window.Render, MainViewWidth/2-300, MainViewHeight/2-300, MainViewWidth/2+300, MainViewHeight/2+300, 255, 0, 0, 180)
	xStart := int32(MainViewWidth/2 - 300)
	yStart := int32(MainViewHeight/2 - 300)
	xEnd := int32(MainViewWidth/2 + 300)
	yEnd := int32(MainViewHeight/2 + 300)

	cells := GetCellInfoMgr().GetAllCellInfo()
	for i, v := range cells {
		x1 := xStart + int32(v.MinX*600/MapWidthWithPixel)
		y1 := yStart + int32(v.MinY*600/MapHeightWithPixel)
		x2 := xStart + int32(v.MaxX*600/MapWidthWithPixel)
		y2 := yStart + int32(v.MaxY*600/MapHeightWithPixel)
		r, g, b := GetRGB(i)
		gfx.BoxRGBA(mv.window.Render, int32(x1), int32(y1), int32(x2), int32(y2), r, g, b, 255)
	}
	gfx.RectangleRGBA(mv.window.Render, xStart, yStart, xEnd, yEnd, 132, 112, 255, 255)

	pos := GetClient().UserView.us.pos
	selfX := xStart + int32(pos.X*600/MapWidthWithPixel)
	selfY := yStart + int32(pos.Z*600/MapHeightWithPixel)
	gfx.FilledCircleRGBA(mv.window.Render, selfX, selfY, 10, 0, 0, 255, 255)
	gfx.CircleRGBA(mv.window.Render, selfX, selfY, 10, 0, 255, 0, 255)
	return true
}
