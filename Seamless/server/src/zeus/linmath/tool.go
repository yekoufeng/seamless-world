package linmath

import (
	"math"
	"math/rand"
	"time"
)

// RandXZ 在XZ平面上半径为r的圆内选取一个随机点
func RandXZ(v Vector3, r float32) Vector3 {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))

	tarR := randSeed.Float64() * float64(r)
	angle := randSeed.Float64() * 2 * math.Pi

	pos := Vector3{}
	pos.Y = 0

	pos.X = float32(math.Cos(angle) * tarR)
	pos.Z = float32(math.Sin(angle) * tarR)

	return v.Add(pos)
}
