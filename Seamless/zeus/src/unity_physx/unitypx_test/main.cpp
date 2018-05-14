#include <iostream>
#include <thread>
#include <random>
#include <chrono>

#include "unitypx.h"

unitypx_sdk_t g_sdk = nullptr;
unitypx_scene_t g_scene = nullptr;

struct ThreadInfo {
	std::unique_ptr<std::thread> thread;
	std::chrono::high_resolution_clock::duration time;
};

const int kRaycastCount = 1000000;
const int kThreadCount = 4;

ThreadInfo g_threads[kThreadCount];

void thread_func(int index, int seed);

int main() {

	g_sdk = unitypx_sdk_create();

	g_scene = unitypx_scene_create(g_sdk, "f:/pxscene");

	if (!g_scene) {
		std::cout << "create scene failed!" << std::endl;
	}
	else {

		std::default_random_engine generator;
		generator.seed((unsigned int)std::chrono::steady_clock::now().time_since_epoch().count());

		for (int i = 0; i < kThreadCount; ++i) {
			g_threads[i].thread.reset(new std::thread(thread_func, i, generator()));
		}

		std::cout << "waiting thread(s)..." << std::endl;

		for (int i = 0; i < kThreadCount; ++i) {
			g_threads[i].thread->join();
		}

		std::cout << "report:" << std::endl;
		std::cout << "Count: " << kRaycastCount << std::endl;
		for (int i = 0; i < kThreadCount; ++i) {
			std::cout << "Thread #" << i + 1 << std::endl;
			std::cout << "\tT:" << std::chrono::duration_cast<std::chrono::duration<double, std::milli>>(g_threads[i].time).count() << " ms" << std::endl;
			std::cout << "\tA:" << std::chrono::duration_cast<std::chrono::duration<double, std::milli>>(g_threads[i].time).count() / kRaycastCount << " ms" << std::endl;
		}

		unitypx_scene_destroy(g_scene);
	}

	unitypx_sdk_destroy(g_sdk);
}

void thread_func(int index, int seed) {
	unitypx_ray_t ray;
	//ray.origin_x = 983;
	ray.origin_y = 1000;
	//ray.origin_z = 970;
	ray.direction_x = 0;
	ray.direction_y = -1;
	ray.direction_z = 0;
	ray.length = 2000;

	unitypx_raycast_result result;

	std::default_random_engine generator;
	std::uniform_real_distribution<float> distribution(0, 2000);
	
	for (int i = 0; i < kRaycastCount; ++i) {
		ray.origin_x = distribution(generator);
		ray.origin_z = distribution(generator);

		auto begin = std::chrono::high_resolution_clock::now();
		unitypx_scene_raycast(g_scene, &ray, 0xFFFFFFFF, &result);
		auto time = std::chrono::high_resolution_clock::now() - begin;

		g_threads[index].time += time;
	}



	//if (unitypx_scene_raycast(g_scene, &ray, 0xFFFFFFFF, &result)) {
	//	std::cout << "hit_distance: " << result.distance << std::endl;
	//	std::cout << "hit_position: " << result.position_x << ", " << result.position_y << ", " << result.position_z << std::endl;
	//	std::cout << "hit_layer: " << result.layer << std::endl;
	//}
}