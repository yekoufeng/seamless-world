package linmath

import (
	"math"
)

// Vector3 代码位置的3D矢量
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// NewVector3 创建一个新的矢量
func NewVector3(x, y, z float32) Vector3 {
	return Vector3{
		x,
		y,
		z,
	}
}

// Vector3_Zero 返回零值
func Vector3_Zero() Vector3 {
	return Vector3{
		0,
		0,
		0,
	}
}

// Vector3_Invalid 返加一个无效的值 ，未赋值之前
func Vector3_Invalid() Vector3 {
	return Vector3{
		math.MaxFloat32,
		math.MaxFloat32,
		math.MaxFloat32,
	}
}

// IsInValid 是否有效
func (v Vector3) IsInValid() bool {
	return v.IsEqual(Vector3_Invalid())
}

// IsEqual 相等
func (v Vector3) IsEqual(r Vector3) bool {
	if v.X-r.X > math.SmallestNonzeroFloat32 ||
		v.X-r.X < -math.SmallestNonzeroFloat32 ||
		v.Y-r.Y > math.SmallestNonzeroFloat32 ||
		v.Y-r.Y < -math.SmallestNonzeroFloat32 ||
		v.Z-r.Z > math.SmallestNonzeroFloat32 ||
		v.Z-r.Z < -math.SmallestNonzeroFloat32 {
		return false
	}

	return true
}

// Add 加
func (v Vector3) Add(o Vector3) Vector3 {
	return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// AddS 加到自己身上
func (v *Vector3) AddS(o Vector3) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

// Sub 减
func (v Vector3) Sub(o Vector3) Vector3 {
	return Vector3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

// SubS 自已身上减
func (v *Vector3) SubS(o Vector3) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
}

// Mul 乘
func (v Vector3) Mul(o float32) Vector3 {
	return Vector3{v.X * o, v.Y * o, v.Z * o}
}

// MulS 自己乘
func (v *Vector3) MulS(o float32) {
	v.X *= o
	v.Y *= o
	v.Z *= o
}

// Cross 叉乘
func (v Vector3) Cross(o Vector3) Vector3 {
	return Vector3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}

// Dot 点乘
func (v Vector3) Dot(o Vector3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Len 获取长度
func (v Vector3) Len() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

func (v *Vector3) Normalize() {
	len := v.Len()

	if len < math.SmallestNonzeroFloat32 {
		return
	}

	v.X = v.X / len
	v.Y = v.Y / len
	v.Z = v.Z / len
}
