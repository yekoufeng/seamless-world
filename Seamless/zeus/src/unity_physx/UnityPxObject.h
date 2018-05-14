#ifndef _UnityPxObject_h_
#define _UnityPxObject_h_

#include "UnityPxSDK.h"

class UnityPxScene;

enum UnityPxObjectType : PxU16 {
	kUnityPxMesh = 1,
	kUnityPxBoxCollider = 2,
	kUnityPxCapsuleCollider = 3,
	kUnityPxMeshCollider = 4,
	kUnityPxTerrainCollider = 5,
};

class UnityPxObject {
public:
	virtual ~UnityPxObject() = default;
	virtual UnityPxObjectType type() = 0;
	virtual void load(PxInputStream &stream) = 0;
	virtual void awake(const UnityPxSDK &sdk, UnityPxScene &scene) = 0;
};

class UnityPxSceneObject : public UnityPxObject {
public:

	~UnityPxSceneObject();

	const PxVec3 & position() const { return m_position; }
	const PxQuat & rotation() const { return m_rotation; }
	PxI32 layer() const { return m_layer; }

	virtual void load(PxInputStream &stream) override;
protected:

	void awake(const UnityPxSDK &sdk, UnityPxScene &scene, PxShape *shape);

	PxVec3 m_position;
	PxQuat m_rotation;
	PxU8 m_layer;

	PxUniquePtr<PxRigidStatic> m_object;
};




#endif