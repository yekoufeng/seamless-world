#ifndef _UnityPxMeshCollider_h_
#define _UnityPxMeshCollider_h_

#include "UnityPxObject.h"

class UnityPxMeshCollider : public UnityPxSceneObject {
public:
	virtual UnityPxObjectType type() override { return kUnityPxMeshCollider; }

	const PxVec3 & scale() const { return m_scale; }

	virtual void load(PxInputStream &stream) override;
	virtual void awake(const UnityPxSDK &sdk, UnityPxScene &scene) override;

protected:
	PxVec3 m_scale;
	PxI32 m_meshIndex;

	PxUniquePtr<PxRigidStatic> m_object;
};

#endif