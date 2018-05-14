package nav

import (
	"errors"
	"zeus/linmath"
)

// Border 边
type Border struct {
	VertA int32
	VertB int32
	NodeA int32
	NodeB int32

	Cost float32
}

func newBorder() *Border {
	return &Border{
		VertA: -1,
		VertB: -1,
		NodeA: -1,
		NodeB: -1,
		Cost:  0,
	}
}

// Node 导航网格结点
type Node struct {
	id      int32
	borders []*Border
	poly    []*linmath.Vector3
	center  linmath.Vector3
	area    int32 //区域信息，约定大于20的值为不可行走区域
}

func newNode() *Node {
	return &Node{
		borders: make([]*Border, 0, 10),
	}
}

func (n *Node) getInteractPoint(pos linmath.Vector3) (linmath.Vector3, int32, error) {

	if len(n.poly) < 3 {
		return linmath.Vector3_Zero(), 0, errors.New("wrong mesh node data ")
	}

	for i := 2; i < len(n.poly); i++ {
		p0 := n.poly[0]
		p1 := n.poly[i-1]
		p2 := n.poly[i]
		poly := []*linmath.Vector3{p0, p1, p2}

		if n.isIncludeInPoly(pos.X, pos.Z, poly) {
			v, err := n.getPlaneInteractPoint(pos, poly)
			return v, n.area, err
		}
	}

	return linmath.Vector3_Zero(), 0, errors.New("no interact point")
}

func (n *Node) getPlaneInteractPoint(pos linmath.Vector3, poly []*linmath.Vector3) (linmath.Vector3, error) {

	lp0 := pos
	ln := linmath.NewVector3(0, -1, 0)

	pp0 := *poly[0]
	pn := poly[1].Sub(*poly[0]).Cross(poly[2].Sub(*poly[0]))

	t := (pn.Dot(pp0) - pn.Dot(lp0)) / pn.Dot(ln)

	if t < 0 {
		return linmath.Vector3_Zero(), errors.New("plane is too high")
	}

	return lp0.Add(ln.Mul(t)), nil
}

func (n *Node) isInclude(x, z float32) bool {
	return n.isIncludeInPoly(x, z, n.poly)
}

func (n *Node) isIncludeInPoly(x, z float32, poly []*linmath.Vector3) bool {

	v := linmath.NewVector3(x, 0, z)

	for i := 0; i < len(poly); i++ {
		v1 := *poly[i]
		v2 := *poly[(i+1)%len(poly)]

		dv := v.Sub(v1)
		de := v2.Sub(v1)

		if n.cross2D(dv, de) < 0 {
			return false
		}
	}

	return true
}

func (n *Node) cross2D(a, b linmath.Vector3) float32 {
	return a.X*b.Z - a.Z*b.X
}
