#include "UnityPxCapsuleCollider.h"

#include "UnityPxScene.h"

void UnityPxCapsuleCollider::load(PxInputStream &stream) {
	UnityPxSceneObject::load(stream);
	stream.read(&m_radius, sizeof(PxF32));
	stream.read(&m_halfHeight, sizeof(PxF32));
}

void UnityPxCapsuleCollider::awake(const UnityPxSDK &sdk, UnityPxScene &scene) {
	auto shape = sdk.physics()->createShape(PxCapsuleGeometry(m_radius, m_halfHeight), *scene.defaultMaterial());
	UnityPxSceneObject::awake(sdk, scene, shape);
	shape->release();
}