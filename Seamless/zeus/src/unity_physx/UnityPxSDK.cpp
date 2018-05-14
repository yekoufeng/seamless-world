#include "UnityPxSDK.h"

UnityPxSDK::~UnityPxSDK() {
	m_pxPhysics.reset();
	m_pxCooking.reset();

	if (m_pxPVD) {
		PxPvdTransport* transport = m_pxPVD->getTransport();
		m_pxPVD.reset();
		transport->release();
	}

	m_pxFoundation.reset();
}

void UnityPxSDK::init(bool pvd) {

	m_pxFoundation.reset(PxCreateFoundation(PX_FOUNDATION_VERSION, m_pxAllocator, m_pxErrorCallback));

	PxCookingParams cookingParams{ PxTolerancesScale() };
	cookingParams.suppressTriangleMeshRemapTable = true;
	cookingParams.midphaseDesc = PxMeshMidPhase::eBVH34;
	cookingParams.midphaseDesc.mBVH34Desc.numTrisPerLeaf = 4;

	m_pxCooking.reset(PxCreateCooking(PX_PHYSICS_VERSION, *m_pxFoundation, cookingParams));

	if (pvd) {
		//m_pxPVD.reset(PxCreatePvd(*m_pxFoundation));
		//PxPvdTransport* transport = PxDefaultPvdSocketTransportCreate("127.0.0.1", 5425, 10);
		//m_pxPVD->connect(*transport, PxPvdInstrumentationFlag::eALL);
	}

	m_pxPhysics.reset(PxCreatePhysics(PX_PHYSICS_VERSION, *m_pxFoundation, PxTolerancesScale(), true, m_pxPVD.get()));

}

static bool intersectRaySphereBasic(const PxVec3& origin, const PxVec3& dir, PxReal length, const PxVec3& center, PxReal radius, PxReal& dist, PxVec3* hit_pos);

bool UnityPxSDK::raycastSphere(const PxVec3 &origin, const PxVec3 &direction, PxF32 length, const PxVec3 &center, PxF32 radius, PxF32 &distance) {
	const PxVec3 x = origin - center;
	PxReal l = PxSqrt(x.dot(x)) - radius - 10.0f;

	//	if(l<0.0f)
	//		l=0.0f;
	l = physx::intrinsics::selectMax(l, 0.0f);

	bool status = intersectRaySphereBasic(origin + l*direction, direction, length - l, center, radius, distance, nullptr);
	if (status)
	{
		//		dist += l/length;
		distance += l;
	}
	return status;
}

static PxReal distancePointSegmentSquaredInternal(const PxVec3& p0, const PxVec3& dir, const PxVec3& point, PxReal* param = NULL);
static PxU32 intersectRayCapsuleInternal(const PxVec3& origin, const PxVec3& dir, const PxVec3& p0, const PxVec3& p1, float radius, PxReal s[2]);

bool UnityPxSDK::raycastCapsule(const PxVec3 &origin, const PxVec3& direction, PxF32 length, const PxVec3 &p0, const PxVec3 &p1, PxF32 radius, PxF32 &distance) {
	// PT: move ray origin close to capsule, to solve accuracy issues.
	// We compute the distance D between the ray origin and the capsule's segment.
	// Then E = D - radius = distance between the ray origin and the capsule.
	// We can move the origin freely along 'dir' up to E units before touching the capsule.
	PxReal l = distancePointSegmentSquaredInternal(p0, p1 - p0, origin);
	l = PxSqrt(l) - radius;

	// PT: if this becomes negative or null, the ray starts inside the capsule and we can early exit
	if (l <= 0.0f)
	{
		distance = 0.0f;
		return true;
	}

	// PT: we remove an arbitrary GU_RAY_SURFACE_OFFSET units to E, to make sure we don't go close to the surface.
	// If we're moving in the direction of the capsule, the origin is now about GU_RAY_SURFACE_OFFSET units from it.
	// If we're moving away from the capsule, the ray won't hit the capsule anyway.
	// If l is smaller than GU_RAY_SURFACE_OFFSET we're close enough, accuracy is good, there is nothing to do.
	if (l > 10.0f)
		l -= 10.0f;
	else
		l = 0.0f;

	// PT: move origin closer to capsule and do the raycast
	PxReal s[2];
	const PxU32 nbHits = intersectRayCapsuleInternal(origin + l*direction, direction, p0, p1, radius, s);
	if (!nbHits)
		return false;

	// PT: keep closest hit only
	if (nbHits == 1)
		distance = s[0];
	else
		distance = (s[0] < s[1]) ? s[0] : s[1];

	// PT: fix distance (smaller than expected after moving ray close to capsule)
	distance += l;
	return true;
}

static bool intersectRaySphereBasic(const PxVec3& origin, const PxVec3& dir, PxReal length, const PxVec3& center, PxReal radius, PxReal& dist, PxVec3* hit_pos)
{
	// get the offset vector
	const PxVec3 offset = center - origin;

	// get the distance along the ray to the center point of the sphere
	const PxReal ray_dist = dir.dot(offset);

	// get the squared distances
	const PxReal off2 = offset.dot(offset);
	const PxReal rad_2 = radius * radius;
	if (off2 <= rad_2)
	{
		// we're in the sphere
		if (hit_pos)
			*hit_pos = origin;
		dist = 0.0f;
		return true;
	}

	if (ray_dist <= 0 || (ray_dist - length) > radius)
	{
		// moving away from object or too far away
		return false;
	}

	// find hit distance squared
	const PxReal d = rad_2 - (off2 - ray_dist * ray_dist);
	if (d<0.0f)
	{
		// ray passes by sphere without hitting
		return false;
	}

	// get the distance along the ray
	dist = ray_dist - PxSqrt(d);
	if (dist > length)
	{
		// hit point beyond length
		return false;
	}

	// sort out the details
	if (hit_pos)
		*hit_pos = origin + dir * dist;
	return true;
}
static PxReal distancePointSegmentSquaredInternal(const PxVec3& p0, const PxVec3& dir, const PxVec3& point, PxReal* param)
{
	PxVec3 diff = point - p0;
	PxReal fT = diff.dot(dir);

	if (fT <= 0.0f)
	{
		fT = 0.0f;
	}
	else
	{
		const PxReal sqrLen = dir.magnitudeSquared();
		if (fT >= sqrLen)
		{
			fT = 1.0f;
			diff -= dir;
		}
		else
		{
			fT /= sqrLen;
			diff -= fT*dir;
		}
	}

	if (param)
		*param = fT;

	return diff.magnitudeSquared();
}
static PxU32 intersectRayCapsuleInternal(const PxVec3& origin, const PxVec3& dir, const PxVec3& p0, const PxVec3& p1, float radius, PxReal s[2])
{
	// set up quadratic Q(t) = a*t^2 + 2*b*t + c
	PxVec3 kW = p1 - p0;
	const float fWLength = kW.magnitude();
	if (fWLength != 0.0f)
		kW /= fWLength;

	// PT: if the capsule is in fact a sphere, switch back to dedicated sphere code.
	// This is not just an optimization, the rest of the code fails otherwise.
	if (fWLength <= 1e-6f)
	{
		const float d0 = (origin - p0).magnitudeSquared();
		const float d1 = (origin - p1).magnitudeSquared();
		const float approxLength = (PxMax(d0, d1) + radius)*2.0f;
		return PxU32(UnityPxSDK::raycastSphere(origin, dir, approxLength, p0, radius, s[0]));
	}

	// generate orthonormal basis
	PxVec3 kU(0.0f);

	if (fWLength > 0.0f)
	{
		PxReal fInvLength;
		if (PxAbs(kW.x) >= PxAbs(kW.y))
		{
			// W.x or W.z is the largest magnitude component, swap them
			fInvLength = PxRecipSqrt(kW.x*kW.x + kW.z*kW.z);
			kU.x = -kW.z*fInvLength;
			kU.y = 0.0f;
			kU.z = kW.x*fInvLength;
		}
		else
		{
			// W.y or W.z is the largest magnitude component, swap them
			fInvLength = PxRecipSqrt(kW.y*kW.y + kW.z*kW.z);
			kU.x = 0.0f;
			kU.y = kW.z*fInvLength;
			kU.z = -kW.y*fInvLength;
		}
	}

	PxVec3 kV = kW.cross(kU);
	kV.normalize();	// PT: fixed november, 24, 2004. This is a bug in Magic.

					// compute intersection

	PxVec3 kD(kU.dot(dir), kV.dot(dir), kW.dot(dir));
	const float fDLength = kD.magnitude();
	const float fInvDLength = fDLength != 0.0f ? 1.0f / fDLength : 0.0f;
	kD *= fInvDLength;

	const PxVec3 kDiff = origin - p0;
	const PxVec3 kP(kU.dot(kDiff), kV.dot(kDiff), kW.dot(kDiff));
	const PxReal fRadiusSqr = radius*radius;

	// Is the velocity parallel to the capsule direction? (or zero)
	if (PxAbs(kD.z) >= 1.0f - PX_EPS_REAL || fDLength < PX_EPS_REAL)
	{
		const float fAxisDir = dir.dot(kW);

		const PxReal fDiscr = fRadiusSqr - kP.x*kP.x - kP.y*kP.y;
		if (fAxisDir < 0 && fDiscr >= 0.0f)
		{
			// Velocity anti-parallel to the capsule direction
			const PxReal fRoot = PxSqrt(fDiscr);
			s[0] = (kP.z + fRoot)*fInvDLength;
			s[1] = -(fWLength - kP.z + fRoot)*fInvDLength;
			return 2;
		}
		else if (fAxisDir > 0 && fDiscr >= 0.0f)
		{
			// Velocity parallel to the capsule direction
			const PxReal fRoot = PxSqrt(fDiscr);
			s[0] = -(kP.z + fRoot)*fInvDLength;
			s[1] = (fWLength - kP.z + fRoot)*fInvDLength;
			return 2;
		}
		else
		{
			// sphere heading wrong direction, or no velocity at all
			return 0;
		}
	}

	// test intersection with infinite cylinder
	PxReal fA = kD.x*kD.x + kD.y*kD.y;
	PxReal fB = kP.x*kD.x + kP.y*kD.y;
	PxReal fC = kP.x*kP.x + kP.y*kP.y - fRadiusSqr;
	PxReal fDiscr = fB*fB - fA*fC;
	if (fDiscr < 0.0f)
	{
		// line does not intersect infinite cylinder
		return 0;
	}

	PxU32 iQuantity = 0;

	if (fDiscr > 0.0f)
	{
		// line intersects infinite cylinder in two places
		const PxReal fRoot = PxSqrt(fDiscr);
		const PxReal fInv = 1.0f / fA;
		PxReal fT = (-fB - fRoot)*fInv;
		PxReal fTmp = kP.z + fT*kD.z;
		const float epsilon = 1e-3f;	// PT: see TA35174
		if (fTmp >= -epsilon && fTmp <= fWLength + epsilon)
			s[iQuantity++] = fT*fInvDLength;

		fT = (-fB + fRoot)*fInv;
		fTmp = kP.z + fT*kD.z;
		if (fTmp >= -epsilon && fTmp <= fWLength + epsilon)
			s[iQuantity++] = fT*fInvDLength;

		if (iQuantity == 2)
		{
			// line intersects capsule wall in two places
			return 2;
		}
	}
	else
	{
		// line is tangent to infinite cylinder
		const PxReal fT = -fB / fA;
		const PxReal fTmp = kP.z + fT*kD.z;
		if (0.0f <= fTmp && fTmp <= fWLength)
		{
			s[0] = fT*fInvDLength;
			return 1;
		}
	}

	// test intersection with bottom hemisphere
	// fA = 1
	fB += kP.z*kD.z;
	fC += kP.z*kP.z;
	fDiscr = fB*fB - fC;
	if (fDiscr > 0.0f)
	{
		const PxReal fRoot = PxSqrt(fDiscr);
		PxReal fT = -fB - fRoot;
		PxReal fTmp = kP.z + fT*kD.z;
		if (fTmp <= 0.0f)
		{
			s[iQuantity++] = fT*fInvDLength;
			if (iQuantity == 2)
				return 2;
		}

		fT = -fB + fRoot;
		fTmp = kP.z + fT*kD.z;
		if (fTmp <= 0.0f)
		{
			s[iQuantity++] = fT*fInvDLength;
			if (iQuantity == 2)
				return 2;
		}
	}
	else if (fDiscr == 0.0f)
	{
		const PxReal fT = -fB;
		const PxReal fTmp = kP.z + fT*kD.z;
		if (fTmp <= 0.0f)
		{
			s[iQuantity++] = fT*fInvDLength;
			if (iQuantity == 2)
				return 2;
		}
	}

	// test intersection with top hemisphere
	// fA = 1
	fB -= kD.z*fWLength;
	fC += fWLength*(fWLength - 2.0f*kP.z);

	fDiscr = fB*fB - fC;
	if (fDiscr > 0.0f)
	{
		const PxReal fRoot = PxSqrt(fDiscr);
		PxReal fT = -fB - fRoot;
		PxReal fTmp = kP.z + fT*kD.z;
		if (fTmp >= fWLength)
		{
			s[iQuantity++] = fT*fInvDLength;
			if (iQuantity == 2)
				return 2;
		}

		fT = -fB + fRoot;
		fTmp = kP.z + fT*kD.z;
		if (fTmp >= fWLength)
		{
			s[iQuantity++] = fT*fInvDLength;
			if (iQuantity == 2)
				return 2;
		}
	}
	else if (fDiscr == 0.0f)
	{
		const PxReal fT = -fB;
		const PxReal fTmp = kP.z + fT*kD.z;
		if (fTmp >= fWLength)
		{
			s[iQuantity++] = fT*fInvDLength;
			if (iQuantity == 2)
				return 2;
		}
	}
	return iQuantity;
}