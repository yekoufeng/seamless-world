#include "UnityPxObject.h"
#include "UnityPxScene.h"

UnityPxSceneObject::~UnityPxSceneObject() {
	if (m_object) {
		m_object->userData = nullptr;
	}
}

void UnityPxSceneObject::load(PxInputStream &stream) {
	stream.read(&m_position, sizeof(PxVec3));
	stream.read(&m_rotation, sizeof(PxQuat));
	stream.read(&m_layer, sizeof(PxU8));
}

void UnityPxSceneObject::awake(const UnityPxSDK &sdk, UnityPxScene &scene, PxShape *shape) {
	PxFilterData data;
	data.word0 = 1 << m_layer;
	shape->setQueryFilterData(data);

	m_object.reset(sdk.physics()->createRigidStatic({ m_position, m_rotation }));
	m_object->attachShape(*shape);
	m_object->userData = this;
	scene.scene()->addActor(*m_object);
}
