#ifndef _unitypx_h_
#define _unitypx_h_

#ifndef UnityPxAPI
#ifdef _WIN32
#	define UnityPxAPI __declspec(dllimport)
#else
#	define UnityPxAPI
#endif
#endif

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {} *unitypx_sdk_t;
typedef struct {} *unitypx_scene_t;
typedef struct {} *unitypx_player_t;

typedef struct {
	float origin_x, origin_y, origin_z;
	float direction_x, direction_y, direction_z;
	float length;
} unitypx_ray_t;

typedef struct {
	float distance;
	float position_x, position_y, position_z;
	int layer;
} unitypx_raycast_result;

typedef struct {
	float center_x, center_y, center_z;
	float radius;
} unitypx_sphere_t;

typedef struct {
	float p0_x, p0_y, p0_z;
	float p1_x, p1_y, p1_z;
	float radius;
} unitypx_capsule_t;

typedef struct {
	float x, y, z;
	float qw, qx, qy, qz;
} unitypx_transform;

UnityPxAPI int unitypx_sphere_raycast(const unitypx_sphere_t *sphere, const unitypx_ray_t *ray, float *distance);

UnityPxAPI int unitypx_capsule_raycast(const unitypx_capsule_t *capsule, const unitypx_ray_t *ray, float *distance);

UnityPxAPI unitypx_sdk_t unitypx_sdk_create();

UnityPxAPI unitypx_scene_t unitypx_scene_create_empty(unitypx_sdk_t sdk);

UnityPxAPI void unitypx_sdk_destroy(unitypx_sdk_t sdk);

UnityPxAPI unitypx_scene_t unitypx_scene_create(unitypx_sdk_t sdk, const char *file);

UnityPxAPI void unitypx_scene_destroy(unitypx_scene_t scene);

UnityPxAPI int unitypx_scene_raycast(unitypx_scene_t scene, const unitypx_ray_t *ray, int mask, unitypx_raycast_result *result);

UnityPxAPI unitypx_player_t unitypx_create_player(unitypx_scene_t scene, float radius, float halfHeight, unitypx_transform* transform);

UnityPxAPI void unitypx_update_player(unitypx_player_t player, unitypx_transform* transform);

#ifdef __cplusplus
}
#endif

#endif