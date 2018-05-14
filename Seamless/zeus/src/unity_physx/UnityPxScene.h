#ifndef _UnityPxScene_h_
#define _UnityPxScene_h_

#include "UnityPxObject.h"

#include <vector>
#include <list>

class UnityPxMesh;

class UnityPxScene {
public:

	~UnityPxScene();

	void init(const UnityPxSDK &sdk);

	bool load(const UnityPxSDK &sdk, PxInputStream &stream);

	static std::unique_ptr<UnityPxObject> createObject(UnityPxObjectType type);

	void addMesh(UnityPxMesh *mesh);
	UnityPxMesh * getMesh(PxI32 index) const;

	const PxUniquePtr<PxScene> & scene() const { return m_pxScene; }
	const PxUniquePtr<PxMaterial> & defaultMaterial() const { return m_pxDefaultMaterial; }

	bool raycast(const PxVec3 &origin, const PxVec3 &direction, PxF32 length, PxU32 mask, PxF32 &distance, PxVec3 &hit, PxI32 &layer) const;

	PxRigidStatic* create_player(float radius, float halfHeight, PxTransform& pos);

protected:



	std::list<std::unique_ptr<UnityPxObject>> m_objects;
	std::vector<UnityPxMesh *> m_meshes;

	PxUniquePtr<PxMaterial> m_pxDefaultMaterial;
	PxUniquePtr<PxDefaultCpuDispatcher> m_pxDispatcher;
	PxUniquePtr<PxScene> m_pxScene;


};

#endif