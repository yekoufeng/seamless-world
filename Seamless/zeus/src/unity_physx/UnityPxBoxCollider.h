#ifndef _UnityPxBoxCollider_h_
#define _UnityPxBoxCollider_h_

#include "UnityPxObject.h"

class UnityPxBoxCollider : public UnityPxSceneObject {
public:
	virtual UnityPxObjectType type() override { return kUnityPxBoxCollider; }

	const PxVec3 & halfExtents() const { return m_halfExtents; }

	virtual void load(PxInputStream &stream) override;
	virtual void awake(const UnityPxSDK &sdk, UnityPxScene &scene) override;

protected:
	PxVec3 m_halfExtents;

	PxUniquePtr<PxRigidStatic> m_object;
};

#endif