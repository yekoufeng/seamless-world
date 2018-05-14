package sdl2

import (
	"errors"
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	IMAGE_SIZE = 10 //像素的大小
)

type MapInfo struct {
	surface *sdl.Surface
	texture *sdl.Texture
}

//type MapArea struct {
//	mapID int32
//	//areaID int32
//	//pos [2]int32
//}

type MapMgr struct {
	winWidth  int32
	winHeight int32
	mapWidth  int32
	mapHeight int32
	renderer  *sdl.Renderer
	mapinfo   []*MapInfo
	//maparea   map[string]*MapArea
}

func GetNewMapMgr(winWidth, winHeight, mapWidth, mapHeight int32, renderer *sdl.Renderer) *MapMgr {
	obj := &MapMgr{
		winWidth:  winWidth,
		winHeight: winHeight,
		renderer:  renderer,
		mapWidth:  mapWidth,
		mapHeight: mapHeight,
		mapinfo:   make([]*MapInfo, 0),
		//maparea:   make(map[string]*MapArea),
	}

	// Todo 待优化
	if winHeight%IMAGE_SIZE != 0 || winWidth%IMAGE_SIZE != 0 {
		panic(errors.New("屏幕设置错误"))
	}
	if mapWidth%IMAGE_SIZE != 0 || mapHeight%IMAGE_SIZE != 0 {
		panic(errors.New("地图设置错误"))
	}

	return obj
}

func (mm *MapMgr) Load(path string, index int) {
	var err error

	mm.mapinfo = mm.mapinfo[:0:0]
	for i := 0; i < index; i += 1 {
		info := &MapInfo{}
		fileName := fmt.Sprintf("%s/%d.png", path, i)
		if info.surface, err = img.Load(fileName); err != nil {
			panic(err)
		}

		if info.texture, err = mm.renderer.CreateTextureFromSurface(info.surface); err != nil {
			panic(err)
		}

		mm.mapinfo = append(mm.mapinfo, info)
	}

	// 随机产生地图 (out of memory)
	// 当前PNG 为10 x 10
	//count := 0
	//for j := int32(0); j < mm.mapHeight/IMAGE_SIZE; j += 1 {
	//	for i := int32(0); i < mm.mapWidth/IMAGE_SIZE; i += 1 {
	//		mapID := rand.Intn(cap(mm.mapinfo))
	//		x := i * IMAGE_SIZE
	//		y := j * IMAGE_SIZE
	//
	//		key := fmt.Sprintf("%d:%d", x, y)
	//		obj := &MapArea{
	//			//areaID: int32(count),
	//			mapID: int32(mapID),
	//		}
	//		//obj.pos[0] = x
	//		//obj.pos[1] = y
	//		mm.maparea[key] = obj
	//
	//		//count += 1
	//	}
	//}
}

func (mm *MapMgr) Destroy() {
	for _, value := range mm.mapinfo {
		value.surface.Free()
		value.texture.Destroy()
	}
	//Todo 切片删除
}

func (mm *MapMgr) Draw(pos [2]int32) {
	var cpRect sdl.Rect
	//cpRect = sdl.Rect{0, 0, IMAGE_SIZE, IMAGE_SIZE}

	//获取顶点
	startPos := pos

	xOffset := pos[0] % IMAGE_SIZE
	zOffset := pos[1] % IMAGE_SIZE

	for j := int32(0); j < mm.winHeight/IMAGE_SIZE+1; j += 1 {
		for i := int32(0); i < mm.winWidth/IMAGE_SIZE+1; i += 1 {

			cpRect = sdl.Rect{X: i*IMAGE_SIZE - xOffset, Y: j*IMAGE_SIZE - zOffset, W: IMAGE_SIZE, H: IMAGE_SIZE}
			var texture *sdl.Texture
			var err error
			texture, err = mm.pos2sdl(startPos[0]+i*IMAGE_SIZE, startPos[1]+j*IMAGE_SIZE)
			if err != nil {
				continue
			}

			mm.renderer.Copy(texture, nil, &cpRect)
		}
	}
}

func (mm *MapMgr) pos2sdl(x, y int32) (texture *sdl.Texture, err error) {
	//取整
	x = x - x%IMAGE_SIZE
	y = y - y%IMAGE_SIZE

	mapId := (x * y) / 8000 % (int32(len(mm.mapinfo)))
	return mm.mapinfo[mapId].texture, nil
}
