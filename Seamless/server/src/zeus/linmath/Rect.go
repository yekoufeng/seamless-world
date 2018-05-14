package linmath

// Rect 定义一个矩形
type Rect struct {
	Xmin float64
	Xmax float64
	Ymin float64
	Ymax float64
}

func (r *Rect) Init(xmin float64, xmax float64, ymin float64, ymax float64) {
	r.Xmin = xmin
	r.Xmax = xmax
	r.Ymin = ymin
	r.Ymax = ymax

}

//IsInOuterBorder 一个点是否位于矩形的外部扩大矩形中
//用于检测一个entity是否移动到其他cell的AOI范围
func (r *Rect) IsInOuterRect(x float64, y float64, length float64) bool {
	Xmin := r.Xmin - length
	Xmax := r.Xmax + length
	Ymin := r.Ymin - length
	Ymax := r.Ymax + length

	if x >= Xmin && x < Xmax && y >= Ymin && y < Ymax {
		return true
	}

	return false
}

//IsInRect 判断一个点是否位于一个内部小矩形中
func (r *Rect) IsInInnerRect(x float64, y float64, length float64) bool {
	Xmin := r.Xmin + length
	Xmax := r.Xmax - length
	Ymin := r.Ymin + length
	Ymax := r.Ymax - length

	if x >= Xmin && x < Xmax && y >= Ymin && y < Ymax {
		return true
	}

	return false
}

//GetArea 获取面积
func (r *Rect) GetArea() float64 {
	area := (r.Xmax - r.Xmin) * (r.Ymax - r.Ymin)

	return area
}
