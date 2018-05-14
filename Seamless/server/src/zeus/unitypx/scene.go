package unitypx

/*
#cgo LDFLAGS: -L./ -lunitypx
#include "unitypx.h"
*/
import "C"

import (
	"fmt"
	"math"
	"zeus/linmath"
)

// Scene 场景静态物理信息
type Scene struct {
	_cScene C.unitypx_scene_t
}

// NewScene 加载新的场景
func NewScene(path string) (*Scene, error) {
	if sdk == nil {
		return nil, fmt.Errorf("加载场景失败, SDK未初始化")
	}
	scene := &Scene{}
	scene._cScene = C.unitypx_scene_create(sdk, C.CString(path+"pxscene"))
	if scene._cScene == nil {
		return nil, fmt.Errorf("加载场景信息失败")
	}
	return scene, nil
}

func NewEmptyScene() *Scene {
	scene := &Scene{}
	scene._cScene = C.unitypx_scene_create_empty(sdk)
	return scene
}

// Raycast 射线检测
func (scene *Scene) Raycast(origin, direction linmath.Vector3, length float32, mask int32) (float32, linmath.Vector3, int32, bool) {
	ray := C.unitypx_ray_t{}
	ray.origin_x = C.float(origin.X)
	ray.origin_y = C.float(origin.Y)
	ray.origin_z = C.float(origin.Z)
	ray.direction_x = C.float(direction.X)
	ray.direction_y = C.float(direction.Y)
	ray.direction_z = C.float(direction.Z)
	ray.length = C.float(length)

	cResult := C.unitypx_raycast_result{}
	cRet := C.unitypx_scene_raycast(scene._cScene, &ray, C.int(mask), &cResult)
	if cRet == 1 {
		pos := linmath.Vector3{}
		pos.X = float32(cResult.position_x)
		pos.Y = float32(cResult.position_y)
		pos.Z = float32(cResult.position_z)

		return float32(cResult.distance), pos, int32(cResult.layer), true
	}

	return 0, linmath.Vector3_Invalid(), 0, false
}

func (scene *Scene) CapsuleRaycast(head, foot linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool) {
	capsule := C.unitypx_capsule_t{}
	capsule.p0_x = C.float(head.X)
	capsule.p0_y = C.float(head.Y)
	capsule.p0_z = C.float(head.Z)
	capsule.p1_x = C.float(foot.X)
	capsule.p1_y = C.float(foot.Y)
	capsule.p1_z = C.float(foot.Z)
	capsule.radius = C.float(radius)

	ray := C.unitypx_ray_t{}
	ray.origin_x = C.float(origin.X)
	ray.origin_y = C.float(origin.Y)
	ray.origin_z = C.float(origin.Z)
	ray.direction_x = C.float(direction.X)
	ray.direction_y = C.float(direction.Y)
	ray.direction_z = C.float(direction.Z)
	ray.length = C.float(length)

	var distance C.float
	cRet := C.unitypx_capsule_raycast(&capsule, &ray, &distance)
	if cRet == 1 {
		return float32(distance), true
	}

	return math.MaxFloat32, false
}

func (scene *Scene) SphereRaycast(center linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool) {
	sphere := C.unitypx_sphere_t{}
	sphere.center_x = C.float(center.X)
	sphere.center_y = C.float(center.Y)
	sphere.center_z = C.float(center.Z)
	sphere.radius = C.float(radius)

	ray := C.unitypx_ray_t{}
	ray.origin_x = C.float(origin.X)
	ray.origin_y = C.float(origin.Y)
	ray.origin_z = C.float(origin.Z)
	ray.direction_x = C.float(direction.X)
	ray.direction_y = C.float(direction.Y)
	ray.direction_z = C.float(direction.Z)
	ray.length = C.float(length)

	var distance C.float
	cRet := C.unitypx_sphere_raycast(&sphere, &ray, &distance)
	if cRet == 1 {
		return float32(distance), true
	}

	return math.MaxFloat32, false
}

func (scene *Scene) CreatePlayer(pos, rota linmath.Vector3) *UnitypxPlayer {
	cpos := C.unitypx_transform{}
	cpos.x = C.float(pos.X)
	cpos.y = C.float(pos.Y)
	cpos.z = C.float(pos.Z)
	cpos.qw = 0
	cpos.qx = C.float(rota.X)
	cpos.qy = C.float(rota.Y)
	cpos.qz = C.float(rota.Z)
	return &UnitypxPlayer{C.unitypx_create_player(scene._cScene, 3.0, 3.0, &cpos)}
}

func (scene *Scene) UpdatePlayer(player *UnitypxPlayer, pos, rota linmath.Vector3) {
	cpos := C.unitypx_transform{}
	cpos.x = C.float(pos.X)
	cpos.y = C.float(pos.Y)
	cpos.z = C.float(pos.Z)
	cpos.qw = 0
	cpos.qx = C.float(rota.X)
	cpos.qy = C.float(rota.Y)
	cpos.qz = C.float(rota.Z)
	C.unitypx_update_player(player.player, &cpos)
}

type UnitypxPlayer struct {
	player C.unitypx_player_t
}
