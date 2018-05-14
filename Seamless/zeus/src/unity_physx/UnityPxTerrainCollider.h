#ifndef _UnityPxTerrainCollider_h_ 
#define _UnityPxTerrainCollider_h_

#include "UnityPxObject.h"

#include <vector>

class UnityPxTerrainCollider : public UnityPxSceneObject {
public:
	
	virtual UnityPxObjectType type() override { return kUnityPxTerrainCollider; }

	virtual void load(PxInputStream &stream) override;
	virtual void awake(const UnityPxSDK &sdk, UnityPxScene &scene) override;

	static PxI16 convertHeightFloatToI16(float value);

protected:
	PxVec3 m_size;
	PxI32 m_heightmapResolution;
	std::vector<PxF32> m_heightmap;

	PxUniquePtr<PxHeightField> m_heightField;

	PxUniquePtr<PxRigidStatic> m_object;
};

#endif
