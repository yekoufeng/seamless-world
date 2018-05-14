package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type ImagePair struct {
	surface *sdl.Surface
	texture *sdl.Texture
}

type AssetsMgr struct {
	images map[string]*ImagePair
}

var _AssetsMgr *AssetsMgr

func GetAssetsMgr() *AssetsMgr {
	if _AssetsMgr == nil {
		_AssetsMgr = &AssetsMgr{
			images: make(map[string]*ImagePair),
		}
	}
	return _AssetsMgr
}

func (assets *AssetsMgr) PreLoad(render *sdl.Renderer) {
	assets.load(render, "assets/map/0.png")
	assets.load(render, "assets/map/1.png")
	assets.load(render, "assets/map/2.png")
	assets.load(render, "assets/map/3.png")

	assets.load(render, "assets/user/1.png")
	assets.load(render, "assets/user/2.png")
	assets.load(render, "assets/user/3.png")
	assets.load(render, "assets/user/4.png")
	assets.load(render, "assets/user/5.png")
	assets.load(render, "assets/user/6.png")
	assets.load(render, "assets/user/7.png")
	assets.load(render, "assets/user/8.png")

	assets.load(render, "assets/cursor.png")

	assets.load(render, "assets/question.jpg")

	assets.load(render, "assets/attack.png")
	// assets.load(render, "assets/circle.png")
}

func (assets *AssetsMgr) load(render *sdl.Renderer, path string) {
	surface, err := img.Load(path)
	if err != nil {
		panic(err)
	}

	texture, err := render.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	assets.images[path] = &ImagePair{surface, texture}
}

func (assets *AssetsMgr) GetImagePair(path string) *ImagePair {
	image, ok := assets.images[path]
	if !ok {
		panic(fmt.Sprintf("resource %v not exist", path))
	}
	return image
}
