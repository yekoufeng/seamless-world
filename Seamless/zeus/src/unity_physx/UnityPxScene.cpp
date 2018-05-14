#include "UnityPxScene.h"
#include "UnityPxMesh.h"
#include "UnityPxBoxCollider.h"
#include "UnityPxCapsuleCollider.h"
#include "UnityPxMeshCollider.h"
#include "UnityPxTerrainCollider.h"


void UnityPxScene::init(const UnityPxSDK &sdk) {
	m_pxDefaultMaterial.reset(sdk.physics()->createMaterial(0.6f, 0.6f, 0));

	PxSceneDesc sceneDesc(sdk.physics()->getTolerancesScale());

	sceneDesc.gravity = PxVec3(0.0f, -9.81f, 0.0f);
	m_pxDispatcher.reset(PxDefaultCpuDispatcherCreate(1));
	sceneDesc.cpuDispatcher = m_pxDispatcher.get();
	sceneDesc.filterShader = PxDefaultSimulationFilterShader;
	m_pxScene.reset(sdk.physics()->createScene(sceneDesc));

	auto client = m_pxScene->getScenePvdClient();
	if (client) {
		client->setScenePvdFlag(PxPvdSceneFlag::eTRANSMIT_CONSTRAINTS, true);
		client->setScenePvdFlag(PxPvdSceneFlag::eTRANSMIT_CONTACTS, true);
		client->setScenePvdFlag(PxPvdSceneFlag::eTRANSMIT_SCENEQUERIES, true);
	}
}


UnityPxScene::~UnityPxScene() {
	m_meshes.clear();
	m_objects.clear();

	m_pxDefaultMaterial.reset();
	m_pxScene.reset();
	m_pxDispatcher.reset();
}


bool UnityPxScene::load(const UnityPxSDK &sdk, PxInputStream &stream) {

	PxU32 header;
	if (stream.read(&header, sizeof(PxU32)) != sizeof(PxU32)) {
		return false;
	}

	if (header != ((PxU32)'P' | ((PxU32)'X' << 8) | ((PxU32)'S' << 16))) {
		return false;
	}

	PxI32 count;
	stream.read(&count, sizeof(PxI32));

	for (PxI32 i = 0; i < count; ++i) {
		UnityPxObjectType type = (UnityPxObjectType)0;
		stream.read(&type, sizeof(UnityPxObjectType));
		auto obj = createObject(type);
		if (!obj) {
			return false;
		}
		obj->load(stream);
		obj->awake(sdk, *this);
		m_objects.push_back(std::move(obj));
	}

	stream.read(&count, sizeof(PxI32));
	for (PxI32 i = 0; i < count; ++i) {
		UnityPxObjectType type = (UnityPxObjectType)0;
		stream.read(&type, sizeof(UnityPxObjectType));
		auto obj = createObject(type);
		if (!obj) {
			return false;
		}
		obj->load(stream);
		obj->awake(sdk, *this);
		m_objects.push_back(std::move(obj));
	}

	m_pxScene->simulate(0.001f);
	m_pxScene->fetchResults(true);

	return true;
}

std::unique_ptr<UnityPxObject> UnityPxScene::createObject(UnityPxObjectType type) {
	switch (type) {
		case kUnityPxMesh:
			return std::unique_ptr<UnityPxObject>(new UnityPxMesh());

		case kUnityPxBoxCollider:
			return std::unique_ptr<UnityPxObject>(new UnityPxBoxCollider());

		case kUnityPxCapsuleCollider:
			return std::unique_ptr<UnityPxObject>(new UnityPxCapsuleCollider());

		case kUnityPxMeshCollider:
			return std::unique_ptr<UnityPxObject>(new UnityPxMeshCollider());

		case kUnityPxTerrainCollider:
			return std::unique_ptr<UnityPxObject>(new UnityPxTerrainCollider());
	}
	return std::unique_ptr<UnityPxObject>();
}

void UnityPxScene::addMesh(UnityPxMesh * mesh) {
	m_meshes.push_back(mesh);
}

UnityPxMesh * UnityPxScene::getMesh(PxI32 index) const {
	if (index < 0 || (size_t)index >= m_meshes.size()) {
		return nullptr;
	}
	return m_meshes[index];
}

bool UnityPxScene::raycast(const PxVec3 & origin, const PxVec3 & direction, PxF32 length, PxU32 mask, PxF32 &distance, PxVec3 &hit, PxI32 &layer) const {
	PxRaycastBuffer buf;
	PxQueryFilterData filterData;
	filterData.data.word0 = mask;
	if (m_pxScene->raycast(origin, direction, length, buf, PxHitFlag::ePOSITION | PxHitFlag::eDISTANCE, filterData)) {
		distance = buf.block.distance;
		hit = buf.block.position;
		layer = (buf.block.actor && buf.block.actor->userData ? static_cast<UnityPxSceneObject *>(buf.block.actor->userData)->layer() : -1);
		return true;
	}
	return false;
}

PxRigidStatic* UnityPxScene::create_player(float radius, float halfHeight, PxTransform& pos) {
	PxCapsuleGeometry geometry(radius, halfHeight);
	auto dynamicShape = m_pxScene->getPhysics().createShape(geometry, *defaultMaterial());

	PxFilterData data;
	data.word0 = (1 << 8);
	dynamicShape->setQueryFilterData(data);

	auto dynamicObject = m_pxScene->getPhysics().createRigidStatic(pos);
	dynamicObject->attachShape(*dynamicShape);
	m_pxScene->addActor(*dynamicObject);
	return dynamicObject;
}
