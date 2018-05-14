
#include "UnityPxSDK.h"
#include "UnityPxScene.h"

#ifndef UnityPxAPI
#ifdef _WIN32
#	define UnityPxAPI __declspec(dllexport)
#else
#	define UnityPxAPI
#endif
#endif

#include <iostream>

extern "C" {

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

	UnityPxAPI int unitypx_sphere_raycast(const unitypx_sphere_t *sphere, const unitypx_ray_t *ray, float *distance) {
		return UnityPxSDK::raycastSphere(
			*reinterpret_cast<const PxVec3 *>(&ray->origin_x),
			*reinterpret_cast<const PxVec3 *>(&ray->direction_x),
			ray->length,
			*reinterpret_cast<const PxVec3 *>(&sphere->center_x),
			sphere->radius,
			*distance
			) ? 1 : 0;
	}

	UnityPxAPI int unitypx_capsule_raycast(const unitypx_capsule_t *capsule, const unitypx_ray_t *ray, float *distance) {
		return UnityPxSDK::raycastCapsule(
			*reinterpret_cast<const PxVec3 *>(&ray->origin_x),
			*reinterpret_cast<const PxVec3 *>(&ray->direction_x),
			ray->length,
			*reinterpret_cast<const PxVec3 *>(&capsule->p0_x),
			*reinterpret_cast<const PxVec3 *>(&capsule->p1_x),
			capsule->radius,
			*distance
			) ? 1 : 0;
	}

	UnityPxAPI unitypx_sdk_t unitypx_sdk_create() {
		auto sdk = new UnityPxSDK();
		sdk->init(false);
		return reinterpret_cast<unitypx_sdk_t>(sdk);
	}

	UnityPxAPI void unitypx_sdk_destroy(unitypx_sdk_t sdk) {
		delete reinterpret_cast<UnityPxSDK *>(sdk);
	}

	UnityPxAPI unitypx_scene_t unitypx_scene_create(unitypx_sdk_t sdk, const char *file) {
		auto _sdk = reinterpret_cast<UnityPxSDK *>(sdk);
		if (!_sdk) {
			return nullptr;
		}
		std::unique_ptr<PxDefaultFileInputData> input{ new PxDefaultFileInputData(file) };
		if (!input->isValid()) {
			return nullptr;
		}

		auto scene = new UnityPxScene();
		scene->init(*_sdk);
		if (!scene->load(*_sdk, *input)) {
			delete scene;
			return nullptr;
		}
		return reinterpret_cast<unitypx_scene_t>(scene);

	}

	UnityPxAPI unitypx_scene_t unitypx_scene_create_empty(unitypx_sdk_t sdk) {
		auto _sdk = reinterpret_cast<UnityPxSDK *>(sdk);
		if (!_sdk) {
			return nullptr;
		}
		auto scene = new UnityPxScene();
		scene->init(*_sdk);

		return reinterpret_cast<unitypx_scene_t>(scene);
	}

	UnityPxAPI void unitypx_scene_destroy(unitypx_scene_t scene) {
		delete reinterpret_cast<UnityPxScene *>(scene);
	}

	UnityPxAPI int unitypx_scene_raycast(unitypx_scene_t scene, const unitypx_ray_t *ray, int mask, unitypx_raycast_result *result) {
		auto _scene = reinterpret_cast<UnityPxScene *>(scene);
		return _scene->raycast(
			*reinterpret_cast<const PxVec3 *>(&ray->origin_x),
			*reinterpret_cast<const PxVec3 *>(&ray->direction_x),
			ray->length,
			mask,
			result->distance,
			*reinterpret_cast<PxVec3 *>(&result->position_x),
			result->layer
			) ? 1 : 0;
	}

	
	UnityPxAPI unitypx_player_t unitypx_create_player(unitypx_scene_t scene, float radius, float halfHeight, unitypx_transform* transform) {
		auto _scene = reinterpret_cast<UnityPxScene *>(scene);
		PxTransform pos({ transform->x, transform->y, transform->z }, { transform->qw, transform->qx, transform->qy, transform->qz });
		auto p = _scene->create_player(radius, halfHeight, pos);
		std::cout << "unitypx_create_player" << std::endl;
		return reinterpret_cast<unitypx_player_t>(p);
	}

	UnityPxAPI void unitypx_update_player(unitypx_player_t player, unitypx_transform* transform) {
		PxRigidStatic* obj = reinterpret_cast<PxRigidStatic*>(player);
		auto pos = obj->getGlobalPose();
		pos.p.x = transform->x;
		pos.p.y = transform->y;
		pos.p.z = transform->z;
		pos.q.w = transform->qw;
		pos.q.x = transform->qx;
		pos.q.y = transform->qy;
		pos.q.z = transform->qz;
		obj->setGlobalPose(pos);
		std::cout << "unitypx_update_player" << std::endl;
	}
}