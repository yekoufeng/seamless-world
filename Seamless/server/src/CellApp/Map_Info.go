package main

import (
	"errors"
	"math"
	"math/rand"
	"zeus/linmath"
)

// MapRange 区域信息
type MapRange struct {
	Typ       int
	ID        int
	CenterPos linmath.Vector3
	Radius    float32
}

func newMapRange() *MapRange {
	return &MapRange{}
}

// GetRandomPos 获取一个随机点
func (r *MapRange) GetRandomPos() linmath.Vector3 {

	rr := rand.Float32() * r.Radius
	ra := rand.Float32() * math.Pi * 2

	x := float32(math.Cos(float64(ra))) * rr
	z := float32(math.Sin(float64(ra))) * rr

	return linmath.NewVector3(x+r.CenterPos.X, r.CenterPos.Y, z+r.CenterPos.Z)
}

// IsContain 是否包含某个位置
func (r *MapRange) IsContain(x, z float32) bool {
	c := linmath.NewVector2(r.CenterPos.X, r.CenterPos.Z)
	v := linmath.NewVector2(x, z)

	return v.Sub(c).Len() < r.Radius
}

// MapRanges 区域列表
type MapRanges struct {
	rangesByTypeID map[int]map[int]*MapRange
	rangesByType   map[int][]*MapRange
}

func newMapRanges() *MapRanges {
	return &MapRanges{
		rangesByTypeID: make(map[int]map[int]*MapRange),
		rangesByType:   make(map[int][]*MapRange),
	}
}

func (mrs *MapRanges) addRange(r *MapRange) {
	ml, ok := mrs.rangesByTypeID[r.Typ]
	if !ok {
		ml = make(map[int]*MapRange)
		mrs.rangesByTypeID[r.Typ] = ml
	}
	ml[r.ID] = r

	ll, ok := mrs.rangesByType[r.Typ]
	if !ok {
		ll = make([]*MapRange, 0, 10)
		mrs.rangesByType[r.Typ] = ll
	}

	mrs.rangesByType[r.Typ] = append(ll, r)
}

// GetRange 获取区域
func (mrs *MapRanges) GetRange(typ, id int) (*MapRange, error) {

	ml, ok := mrs.rangesByTypeID[typ]
	if !ok {
		return nil, errors.New("no type")
	}

	r, ok := ml[id]
	if !ok {
		return nil, errors.New("no id")
	}

	return r, nil
}

// GetRangeList 获取区域列表
func (mrs *MapRanges) GetRangeList(typ int) ([]*MapRange, error) {

	ll, ok := mrs.rangesByType[typ]
	if !ok {
		return nil, errors.New("no type")
	}

	return ll, nil
}

// GetRangeByPos 根据点获取一个区域
func (mrs *MapRanges) GetRangeByPos(x, z float32, typ int) (*MapRange, error) {

	ll, err := mrs.GetRangeList(typ)
	if err != nil {
		return nil, err
	}

	for _, r := range ll {
		if r.IsContain(x, z) {
			return r, nil
		}
	}
	return nil, errors.New("no range")
}

// MapHeightMap 高度图
type MapHeightMap struct {
	Width  float32
	Height float32
	OrgX   float32
	OrgZ   float32
	Res    int
	Data   []float32
}

func newMapHeightMap(w, h, ox, oz float32, r uint32) *MapHeightMap {
	return &MapHeightMap{
		Width:  w,
		Height: h,
		OrgX:   ox,
		OrgZ:   oz,
		Res:    int(r),
		Data:   make([]float32, r*r),
	}
}

// GetHeight 获取高度图
func (hm *MapHeightMap) GetHeight(x, z float32) (float32, error) {

	x -= hm.OrgX
	z -= hm.OrgZ

	if x < 0 || x > hm.Width || z < 0 || z > hm.Height {
		return 0, errors.New("wrong cordinate ")
	}

	xr := x / float32(hm.Width) * float32(hm.Res-1)
	zr := z / float32(hm.Height) * float32(hm.Res-1)

	x1 := float32(math.Floor(float64(xr)))
	x2 := float32(math.Ceil(float64(xr)))
	z1 := float32(math.Floor(float64(zr)))
	z2 := float32(math.Ceil(float64(zr)))

	wx := xr - x1
	wz := zr - z1

	x1z1 := hm.getHeight(int(z1), int(x1))
	x2z2 := hm.getHeight(int(z2), int(x2))
	x1z2 := hm.getHeight(int(z2), int(x1))
	x2z1 := hm.getHeight(int(z1), int(x2))

	cx1 := x1z1 + (x2z1-x1z1)*wx
	cx2 := x1z2 + (x2z2-x1z2)*wx
	fh := cx1 + (cx2-cx1)*wz

	return fh, nil
}

// IsWater 是否水域，水平面为0
func (hm *MapHeightMap) IsWater(x, z float32) (bool, error) {
	h, err := hm.GetHeight(x, z)
	if err != nil {
		return false, err
	}

	return h < 5.5, nil
}

func (hm *MapHeightMap) getHeight(row, col int) float32 {
	return hm.Data[row*hm.Res+col]
}
