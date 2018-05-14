#ifndef _UnityPxSDK_h_
#define _UnityPxSDK_h_

#include "PxPhysicsAPI.h"
#include <memory>

using namespace physx;

template <class T>
struct PxDelete {
	void operator ()(T *obj) const { obj->release(); }
};

template <class T>
struct PxUniquePtr : std::unique_ptr<T, PxDelete<T>> {
	using unique_ptr::unique_ptr;
};


class UnityPxSDK {
public:

	~UnityPxSDK();

	void init(bool pvd);

	const PxUniquePtr<PxCooking> & cooking() const { return m_pxCooking; }
	const PxUniquePtr<PxPhysics> & physics() const { return m_pxPhysics; }

	static bool raycastSphere(const PxVec3 &origin, const PxVec3 &direction, PxF32 length, const PxVec3 &center, PxF32 radius, PxF32 &distance);
	static bool raycastCapsule(const PxVec3 &origin, const PxVec3 &direction, PxF32 length, const PxVec3 &p0, const PxVec3 &p1, PxF32 radius, PxF32 &distance);

protected:
	PxDefaultAllocator		m_pxAllocator;
	PxDefaultErrorCallback	m_pxErrorCallback;

	PxUniquePtr<PxPhysics> m_pxPhysics;
	PxUniquePtr<PxCooking>  m_pxCooking;
	PxUniquePtr<PxPvd> m_pxPVD;
	PxUniquePtr<PxFoundation> m_pxFoundation;
};


#endif