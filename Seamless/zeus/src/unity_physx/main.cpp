#include "UnityPxScene.h"

#include <iostream>

int main() {

	UnityPxSDK sdk;

	sdk.init(false);

	{

		UnityPxScene scene1, scene2;
		scene1.init(sdk);
		scene2.init(sdk);



		{
			std::unique_ptr<PxDefaultFileInputData> file{ new PxDefaultFileInputData("f:\\pxscene1") };
			scene2.load(sdk, *file);
		}




		{
			std::unique_ptr<PxDefaultFileInputData> file{ new PxDefaultFileInputData("f:\\pxscene") };
			scene1.load(sdk, *file);
		}


		PxRaycastBuffer buf;
		PxQueryFilterData filterData;
		filterData.data.word0 = 1 << 12;
		if (scene1.scene()->raycast({ 983, 500, 970 }, { 0, -1, 0 }, 1000, buf, PxHitFlag::ePOSITION, filterData)) {
			std::cout << "hit_position: " << buf.block.position.x << ", " << buf.block.position.y << ", " << buf.block.position.z << std::endl;
			std::cout << "hit_layer: " << (buf.block.actor->userData ? static_cast<UnityPxSceneObject *>(buf.block.actor->userData)->layer() : -1) << std::endl;
		}



		std::cout << "scene1: " << scene1.scene()->getNbActors(PxActorTypeFlag::eRIGID_STATIC) << std::endl;
		std::cout << "scene2: " << scene2.scene()->getNbActors(PxActorTypeFlag::eRIGID_STATIC) << std::endl;

	}

	std::cout << "ok" << std::endl;
	std::cin.get();
}