package main

import (
	"fmt"
	"sync"
	"zeus/nav"
	"zeus/unitypx"

	log "github.com/cihub/seelog"
)

// Maps 地图信息集合
type Maps struct {
	sync.Mutex
	loadedMap  map[string]*MapTask
	pendingMap map[string]*MapTask
}

var mapsInst *Maps

func getMapsInst() *Maps {
	if mapsInst == nil {
		mapsInst = &Maps{
			loadedMap:  make(map[string]*MapTask),
			pendingMap: make(map[string]*MapTask),
		}
	}

	return mapsInst
}

func (ms *Maps) getMap(mapName string) *MapTask {

	ms.Lock()
	defer ms.Unlock()

	task, ok := ms.loadedMap[mapName]

	if !ok {
		task, ok = ms.pendingMap[mapName]
		if !ok {
			task = newMapTask(mapName)
			ms.pendingMap[mapName] = task
			go ms.loadMap(task)
		}
	}

	return task
}

func (ms *Maps) loadMap(task *MapTask) {
	path := fmt.Sprint("../res/space/", task.name, "/")
	mesh, pathFinder, err := ms.loadNavMesh(path)

	loadSucceed := true
	/*
		if err != nil {
			log.Error("load navmesh failed ", path)
			loadSucceed = false
		}
	*/

	ranges, err := ms.loadRange(path)
	if err != nil {
		log.Error("load range info failed ", path)
		loadSucceed = false
	}

	scene, err := unitypx.NewScene(path)
	if err != nil {
		log.Error("load scene info failed ", path)
		loadSucceed = false
	}

	mapInfo := newMap()
	mapInfo.mesh = mesh
	mapInfo.pathFinder = pathFinder
	mapInfo.ranges = ranges
	mapInfo.scene = scene

	task.info = mapInfo

	ms.Lock()

	delete(ms.pendingMap, task.name)
	ms.loadedMap[task.name] = task

	ms.Unlock()

	task.c <- loadSucceed
}

// MapTask 地图任务
type MapTask struct {
	name string
	info *Map
	c    chan bool
}

func newMapTask(name string) *MapTask {
	return &MapTask{
		name,
		nil,
		make(chan bool, 1),
	}
}

// Map 单独的地图信息
type Map struct {
	mesh       *nav.Mesh
	pathFinder *nav.MeshPathFinder
	ranges     *MapRanges
	scene      *unitypx.Scene
}

func newMap() *Map {
	return &Map{}
}
