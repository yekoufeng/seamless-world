#include "UnityPxMeshCollider.h"

#include "UnityPxScene.h"
#include "UnityPxMesh.h"

void UnityPxMeshCollider::load(PxInputStream &stream) {
	UnityPxSceneObject::load(stream);
	stream.read(&m_scale, sizeof(PxVec3));
	stream.read(&m_meshIndex, sizeof(PxI32));
}

void UnityPxMeshCollider::awake(const UnityPxSDK &sdk, UnityPxScene &scene) {
	auto mesh = scene.getMesh(m_meshIndex);
	if (mesh) {
		auto shape = sdk.physics()->createShape(PxTriangleMeshGeometry(mesh->mesh().get(), { m_scale }), *scene.defaultMaterial());
		UnityPxSceneObject::awake(sdk, scene, shape);
		shape->release();
	}
}