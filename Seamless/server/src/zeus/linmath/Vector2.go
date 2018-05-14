package linmath

import "math"

// Vector2 代表位置的2D矢量
type Vector2 struct {
	X float32
	Y float32
}

// NewVector2 创建一个新的Vector2
func NewVector2(x, y float32) Vector2 {
	return Vector2{
		x,
		y,
	}
}

// Vector2_Zero 返回零值
func Vector2_Zero() Vector2 {
	return Vector2{
		0,
		0,
	}
}

// Vector2_Invalid 返加一个无效的值 ，未赋值之前
func Vector2_Invalid() Vector2 {
	return Vector2{
		math.MaxFloat32,
		math.MaxFloat32,
	}
}

// IsInValid 是否有效
func (v Vector2) IsInValid() bool {
	return v.IsEqual(Vector2_Invalid())
}

// IsEqual 相等
func (v Vector2) IsEqual(r Vector2) bool {
	return v.X == r.X && v.Y == r.Y
}

// Add 加
func (v Vector2) Add(o Vector2) Vector2 {
	return Vector2{v.X + o.X, v.Y + o.Y}
}

// AddS 加到自己身上
func (v *Vector2) AddS(o Vector2) {
	v.X += o.X
	v.Y += o.Y
}

// Sub 减
func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{v.X - o.X, v.Y - o.Y}
}

// SubS 自已身上减
func (v *Vector2) SubS(o Vector2) {
	v.X -= o.X
	v.Y -= o.Y
}

// Mul 乘
func (v Vector2) Mul(m float32) Vector2 {
	return Vector2{
		v.X * m,
		v.Y * m,
	}
}

// Dot 点乘
func (v Vector2) Dot(o Vector2) float32 {
	return v.X*o.X + v.Y*o.Y
}

// Len 获取长度
func (v Vector2) Len() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

// Cross 叉乘
func (v Vector2) Cross(o Vector2) float32 {
	return v.X*o.Y - v.Y*o.X
}

func (v *Vector2) Normalize() {
	len := v.Len()
	if len < math.SmallestNonzeroFloat32 {
		return
	}

	v.X = v.X / len
	v.Y = v.Y / len
}
