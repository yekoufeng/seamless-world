#ifndef _UnityPxCapsuleCollider_h_
#define _UnityPxCapsuleCollider_h_

#include "UnityPxObject.h"

class UnityPxCapsuleCollider : public UnityPxSceneObject {
public:
	virtual UnityPxObjectType type() override { return kUnityPxCapsuleCollider; }

	const PxF32 radius() const { return m_radius; }
	const PxF32 halfHeight() const { return m_halfHeight; }

	virtual void load(PxInputStream &stream) override;
	virtual void awake(const UnityPxSDK &sdk, UnityPxScene &scene) override;
protected:
	PxF32 m_radius;
	PxF32 m_halfHeight;

	PxUniquePtr<PxRigidStatic> m_object;
};

#endif
