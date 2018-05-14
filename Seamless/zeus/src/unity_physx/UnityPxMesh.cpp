#include "UnityPxMesh.h"
#include "UnityPxScene.h"

void UnityPxMesh::load(PxInputStream &stream) {
	PxI32 count;
	stream.read(&count, sizeof(PxI32));
	m_vertices.assign(count, PxVec3());
	stream.read(m_vertices.data(), sizeof(PxVec3) * count);
	stream.read(&count, sizeof(PxI32));
	m_indices.assign(count, 0);
	stream.read(m_indices.data(), sizeof(PxU16) * count);
}

void UnityPxMesh::awake(const UnityPxSDK &sdk, UnityPxScene &scene) {
	scene.addMesh(this);

	PxTriangleMeshDesc desc;
	desc.points.count = (PxU32)m_vertices.size();
	desc.points.data = m_vertices.data();
	desc.points.stride = sizeof(PxVec3);
	desc.triangles.count = (PxU32)m_indices.size() / 3;
	desc.triangles.data = m_indices.data();
	desc.triangles.stride = 3 * sizeof(PxU16);
	desc.flags = PxMeshFlag::e16_BIT_INDICES;

	auto mesh = sdk.cooking()->createTriangleMesh(desc, sdk.physics()->getPhysicsInsertionCallback());

	m_mesh.reset(mesh);

	m_vertices.clear();
	m_indices.clear();
}