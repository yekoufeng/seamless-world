#ifndef _UnityPxMesh_h_
#define _UnityPxMesh_h_

#include "UnityPxObject.h"

#include <vector>

class UnityPxMesh : public UnityPxObject {
public:
	virtual UnityPxObjectType type() override { return kUnityPxMesh; }

	virtual void load(PxInputStream &stream) override;
	virtual void awake(const UnityPxSDK &sdk, UnityPxScene &scene) override;

	const PxUniquePtr<PxTriangleMesh> & mesh() const { return m_mesh; }

protected:
	std::vector<PxVec3> m_vertices;
	std::vector<PxU16> m_indices;
	PxUniquePtr<PxTriangleMesh> m_mesh;
};

#endif