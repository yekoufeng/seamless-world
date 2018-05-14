#include "UnityPxBoxCollider.h"

#include "UnityPxScene.h"

void UnityPxBoxCollider::load(PxInputStream &stream) {
	UnityPxSceneObject::load(stream);
	stream.read(&m_halfExtents, sizeof(PxVec3));
}

void UnityPxBoxCollider::awake(const UnityPxSDK &sdk, UnityPxScene &scene) {
	auto shape = sdk.physics()->createShape(PxBoxGeometry(m_halfExtents), *scene.defaultMaterial());
	UnityPxSceneObject::awake(sdk, scene, shape);
	shape->release();
}
