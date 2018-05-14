package nav

import (
	"errors"
	"math"
	"zeus/linmath"
)

//MeshPathFinder 寻路
type MeshPathFinder struct {
	mesh *Mesh
}

// NewMeshPathFinder 创建一个建路器
func NewMeshPathFinder(mesh *Mesh) *MeshPathFinder {
	return &MeshPathFinder{
		mesh: mesh,
	}
}

// FindPath 寻路
func (p *MeshPathFinder) FindPath(srcPos, destPos linmath.Vector3) ([]linmath.Vector3, error) {

	var err error
	var srcNode, destNode *Node

	srcPos, srcNode, err = p.mesh.GetPointAndNode(srcPos)
	if err != nil {
		return nil, err
	}

	if srcNode.area >= 20 {
		return nil, errors.New("srcPos isn't walkable")
	}

	destPos, destNode, err = p.mesh.GetPointAndNode(destPos)
	if err != nil {
		return nil, err
	}

	if destNode.area >= 20 {
		return nil, errors.New("destPos isn't walkable")
	}

	if srcNode == destNode {
		path := make([]linmath.Vector3, 0, 2)
		path = append(path, srcPos)
		path = append(path, destPos)
		return path, nil
	}

	path, err := p.astarFind(srcNode, destNode)

	/*
		fmt.Println("A* 寻路点 开始")

		fmt.Println(path)
		for _, v := range path {
			fmt.Println(*v)
		}

		if err != nil {
			return nil, err
		}

		fmt.Println("A* 寻路点 结束")
	*/

	return p.findWayPoint(srcPos, destPos, path)
}

func (p *MeshPathFinder) astarFind(srcNode, destNode *Node) ([]*Node, error) {
	// 此函数考虑多线程调用，所以每次寻路都使用一个新的astar对象
	star := newAStar(p.mesh)
	return star.FindPath(srcNode, destNode)
}

func (p *MeshPathFinder) findWayPoint(srcPos, destPos linmath.Vector3, path []*Node) ([]linmath.Vector3, error) {
	// 2D方向查找路点

	wayPoints := make([]linmath.Vector3, 0, 10)

	curPos := srcPos
	wayPoints = append(wayPoints, curPos)

	leftEdge := linmath.Vector2_Invalid()
	rightEdge := linmath.Vector2_Invalid()

	leftVert := linmath.Vector3_Invalid()
	rightVert := linmath.Vector3_Invalid()

	var lastCenterPos, leftCorner, rightCorner linmath.Vector3

	for i := 0; i < len(path)-1; i++ {
		curNode := path[i]
		nextNode := path[i+1]

		leftCorner, rightCorner = p.mesh.getCorner(curNode, nextNode)

		lastCenterPos = leftCorner.Add(rightCorner).Mul(0.5)

		//fmt.Println("find node   ", curNode.id, nextNode.id, "   ver =  ", leftCorner, rightCorner)

		newLeftEdge := p.getRay(curPos, leftCorner)
		newRightEdge := p.getRay(curPos, rightCorner)

		if leftEdge.IsInValid() {
			leftEdge = newLeftEdge
			rightEdge = newRightEdge

			leftVert = leftCorner
			rightVert = rightCorner
			continue
		}

		nlol := newLeftEdge.Cross(leftEdge)
		nlor := newLeftEdge.Cross(rightEdge)
		nrol := newRightEdge.Cross(leftEdge)
		nror := newRightEdge.Cross(rightEdge)

		if nlol > 0 && nlor < 0 { //新的左射线在范围内

			//fmt.Println("replace left ")

			leftEdge = newLeftEdge
			leftVert = leftCorner
		}
		if nrol > 0 && nror < 0 { // 新的右射线在范围内

			//fmt.Println("replace right")

			rightEdge = newRightEdge
			rightVert = rightCorner
		}
		if nlol < 0 && nrol < 0 { //新的射线都在左拐角外

			//fmt.Println("add node ", curNode.id, "  ---  ", nextNode.id, "  LeftCorner ", leftVert)

			curPos = leftVert
			wayPoints = append(wayPoints, curPos)
			leftEdge = p.getRay(curPos, leftCorner)
			rightEdge = p.getRay(curPos, rightCorner)

			leftVert = leftCorner
			rightVert = rightCorner
		} else if nlor > 0 && nror > 0 { //新的射线都在右拐角外

			//fmt.Println("add node ", curNode.id, "  ---  ", nextNode.id, "  RightCorner ", rightVert)

			curPos = rightVert
			wayPoints = append(wayPoints, curPos)
			leftEdge = p.getRay(curPos, leftCorner)
			rightEdge = p.getRay(curPos, rightCorner)

			leftVert = leftCorner
			rightVert = rightCorner
		}
	}

	ray := p.getRay(curPos, destPos)

	if ray.Cross(leftEdge) < 0 || ray.Cross(rightEdge) > 0 {

		iv, e := p.calInterPoint(curPos, leftEdge.Add(rightEdge).Mul(0.5), leftCorner, rightCorner)
		if e == nil {
			lastCenterPos = iv
			//fmt.Println("add last centerPos ", lastCenterPos)
		}

		wayPoints = append(wayPoints, lastCenterPos)
	}

	wayPoints = append(wayPoints, destPos)

	return wayPoints, nil
}

func (p *MeshPathFinder) getRay(startPos, endPos linmath.Vector3) linmath.Vector2 {
	var ray linmath.Vector2

	ray.X = endPos.X - startPos.X
	ray.Y = endPos.Z - startPos.Z

	return ray
}

func (p *MeshPathFinder) calInterPoint(origin linmath.Vector3, normal linmath.Vector2, p0, p1 linmath.Vector3) (linmath.Vector3, error) {

	aStartPoint := linmath.NewVector2(origin.X, origin.Z)
	aEndPoint := aStartPoint.Add(normal)

	bStartPoint := linmath.NewVector2(p0.X, p0.Z)
	bEndPoint := linmath.NewVector2(p1.X, p1.Z)

	ip, err := p.getInterHPoint(aStartPoint, aEndPoint, bStartPoint, bEndPoint)

	return linmath.NewVector3(ip.X, p0.Y, ip.Y), err
}

func (p *MeshPathFinder) getInterHPoint(aStartPoint, aEndPoint, bStartPoint, bEndPoint linmath.Vector2) (linmath.Vector2, error) {
	//行列式求两条直线交点

	D := (aEndPoint.X-aStartPoint.X)*(bStartPoint.Y-bEndPoint.Y) - (bEndPoint.X-bStartPoint.X)*(aStartPoint.Y-aEndPoint.Y)
	D1 := (bStartPoint.Y*bEndPoint.X-bStartPoint.X*bEndPoint.Y)*(aEndPoint.X-aStartPoint.X) - (aStartPoint.Y*aEndPoint.X-aStartPoint.X*aEndPoint.Y)*(bEndPoint.X-bStartPoint.X)
	D2 := (aStartPoint.Y*aEndPoint.X-aStartPoint.X*aEndPoint.Y)*(bStartPoint.Y-bEndPoint.Y) - (bStartPoint.Y*bEndPoint.X-bStartPoint.X*bEndPoint.Y)*(aStartPoint.Y-aEndPoint.Y)

	if math.Abs((float64)(D)) < math.SmallestNonzeroFloat64 {
		return linmath.Vector2_Zero(), errors.New("wrong intersection point")
	}

	return linmath.NewVector2(D1/D, D2/D), nil
}
