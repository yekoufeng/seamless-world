package main

import (
	"math"
	"protoMsg"
	"sync"
	"sync/atomic"
	"zeus/linmath"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	KeyMask_W = 0x1
	KeyMask_S = 0x2
	KeyMask_A = 0x4
	KeyMask_D = 0x8
)

const (
	MoveTp_None     = 0
	MoveTp_Keyboard = 1
	MoveTp_Mouse    = 2
)

type UserState struct {
	pos          linmath.Vector3
	rota         linmath.Vector3
	speed        float64
	lastMoveTime int64

	mask int32

	mtx *sync.Mutex

	keyMask  int
	moveType int32

	clickPos linmath.Vector3
}

func (us *UserState) CloneState(ret *UserState) {
	us.mtx.Lock()
	defer us.mtx.Unlock()

	*ret = *us
	ret.mtx = nil
}

func (us *UserState) SetPos(pos linmath.Vector3) {
	us.mtx.Lock()
	defer us.mtx.Unlock()

	us.pos = pos
	us.AddMask((0x1 | 0x2 | 0x4))
}

func (us *UserState) SetRota(rota linmath.Vector3) {
	us.mtx.Lock()
	defer us.mtx.Unlock()

	us.rota = rota
	us.AddMask((0x8 | 0x10 | 0x20))
}
func (us *UserState) SyncMoveToCell() {
	us.mtx.Lock()
	defer us.mtx.Unlock()

	if atomic.LoadInt32(&us.mask) == 0 {
		return
	}

	stoped := atomic.LoadInt32(&us.moveType) == MoveTp_None
	msg := &protoMsg.MoveReq{
		Pos: &protoMsg.Vector3{
			X: us.pos.X / MapRate,
			Z: us.pos.Z / MapRate,
		},
		Rota: &protoMsg.Vector3{
			X: us.rota.X,
			Z: us.rota.Z,
		},
		Stoped: stoped,
	}

	GetClient().SendMsgToCell(msg)
	us.ResetMask()
}

func (us *UserState) GetAngle() float64 {
	return GetVecAngle(&us.rota)
}

func (us *UserState) UpdateMove() {
	moveType := atomic.LoadInt32(&us.moveType)
	if moveType == MoveTp_Keyboard {
		us.MoveCtrlByKeyboard()
	} else if moveType == MoveTp_Mouse {
		us.MoveCtrlByMouse()
	}
}

func (us *UserState) ViewPosToMapPos(pos linmath.Vector3) linmath.Vector3 {
	rect := us.GetViewMapRect()
	pos.X += float32(rect.X)
	pos.Z += float32(rect.Y)
	return pos
}

func (us *UserState) GetViewMapRect() sdl.Rect {
	rect := sdl.Rect{
		W: MainViewWidth,
		H: MainViewHeight,
	}
	x, z := us.pos.X, us.pos.Z
	if x > MainViewWidth/2 {
		if x > MapWidthWithPixel-MainViewWidth/2 {
			rect.X = MapWidthWithPixel - MainViewWidth
		} else {
			rect.X = int32(x) - MainViewWidth/2
		}
	}
	if z > MainViewHeight/2 {
		if z > MapHeightWithPixel-MainViewHeight/2 {
			rect.Y = MapHeightWithPixel - MainViewHeight
		} else {
			rect.Y = int32(z) - MainViewHeight/2
		}
	}
	return rect
}

func (us *UserState) GetUserViewPos() (float32, float32) {
	//这都都是用主角的坐标算出地图区域
	rect := GetClient().UserView.us.GetViewMapRect()
	return us.pos.X - float32(rect.X), us.pos.Z - float32(rect.Y)
}

func (us *UserState) VerifyPos() {
	if us.pos.X < 30 {
		us.pos.X = 30
	} else if us.pos.X > MapWidthWithPixel-32 {
		us.pos.X = MapWidthWithPixel - 32
	}
	if us.pos.Z < 40 {
		us.pos.Z = 40
	} else if us.pos.Z > MapHeightWithPixel-15 {
		us.pos.Z = MapHeightWithPixel - 15
	}
}

func (us *UserState) AddMask(m int32) {
	mask := atomic.LoadInt32(&us.mask)
	mask |= m
	atomic.StoreInt32(&us.mask, mask)
}

func (us *UserState) ResetMask() {
	atomic.StoreInt32(&us.mask, 0)
}

func GetVecAngle(v *linmath.Vector3) float64 {
	x, z := v.X, v.Z
	l := math.Sqrt(float64(x*x + z*z))
	if l == 0 {
		return 0
	}
	ret := math.Acos(float64(z) / l)
	if x < 0 {
		ret = math.Pi*2 - ret
	}
	for ret < 0 {
		ret += math.Pi * 2
	}
	for ret > math.Pi*2 {
		ret -= math.Pi * 2
	}
	return ret
}
