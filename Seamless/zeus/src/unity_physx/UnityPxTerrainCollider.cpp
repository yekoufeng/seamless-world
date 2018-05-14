#include "UnityPxTerrainCollider.h"
#include "UnityPxScene.h"
#include <cmath>

void UnityPxTerrainCollider::load(PxInputStream &stream) {
	UnityPxSceneObject::load(stream);
	stream.read(&m_size, sizeof(PxVec3));
	stream.read(&m_heightmapResolution, sizeof(PxI32));
	m_heightmap.assign(m_heightmapResolution * m_heightmapResolution, 0);
	stream.read(m_heightmap.data(), (PxU32)m_heightmap.size() * sizeof(PxF32));
}

void UnityPxTerrainCollider::awake(const UnityPxSDK &sdk, UnityPxScene &scene) {


	auto samples = new PxHeightFieldSample[m_heightmapResolution * m_heightmapResolution];

	int materialIndex = 0;

	for (int y = 0; y < m_heightmapResolution; ++y) {
		int rowbase = y * m_heightmapResolution;
		for (int x = 0; x < m_heightmapResolution; ++x) {
			int heightIndex = rowbase + x;
			auto &sample = samples[x * m_heightmapResolution + y];
			sample.height = convertHeightFloatToI16(m_heightmap[heightIndex]);
			sample.materialIndex0 = sample.materialIndex1 = materialIndex;
			sample.setTessFlag();
		}
	}

	PxHeightFieldDesc desc;
	desc.convexEdgeThreshold = 4;
	desc.nbRows = m_heightmapResolution;
	desc.nbColumns = m_heightmapResolution;
	desc.samples.data = samples;
	desc.samples.stride = sizeof(PxHeightFieldSample);

	auto heightField = sdk.cooking()->createHeightField(desc, sdk.physics()->getPhysicsInsertionCallback());

	m_heightField.reset(heightField);

	delete[]samples;
	m_heightmap.clear();

	auto shape = sdk.physics()->createShape(PxHeightFieldGeometry(m_heightField.get(), {}, m_size.y / 32766, m_size.x / (m_heightmapResolution - 1), m_size.z / (m_heightmapResolution - 1)), *scene.defaultMaterial());
	UnityPxSceneObject::awake(sdk, scene, shape);
	shape->release();
}


PxI16 UnityPxTerrainCollider::convertHeightFloatToI16(float value) {
	return (PxI16)PxClamp((int)roundf(value * 32766), 0, 32766);
}