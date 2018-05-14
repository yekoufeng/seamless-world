package nav

import (
	"errors"
	"zeus/linmath"

	log "github.com/cihub/seelog"
)

// Mesh 导航网格
type Mesh struct {
	Verts          []linmath.Vector3
	Nodes          []*Node
	BordersByNodes map[int64]*Border
	BordersByVerts map[int64]*Border
	NodePosMap     map[int32]map[int32][]int32
}

func newMesh() *Mesh {
	return &Mesh{
		Verts:          make([]linmath.Vector3, 0, 1000),
		Nodes:          make([]*Node, 0, 1000),
		BordersByNodes: make(map[int64]*Border),
		BordersByVerts: make(map[int64]*Border),
		NodePosMap:     make(map[int32]map[int32][]int32),
	}
}

// GetPoint 获得当前点往下相交的navmesh的交点以及区域信息
func (m *Mesh) GetPoint(pos linmath.Vector3) (linmath.Vector3, int32, error) {

	ns, err := m.getNodeIDByPos(pos.X, pos.Z)
	if err != nil {
		return linmath.Vector3_Zero(), 0, err
	}

	for _, n := range ns {
		v, a, err := n.getInteractPoint(pos)
		if err == nil {
			return v, a, nil
		}
	}

	return linmath.Vector3_Zero(), 0, errors.New("no interact point")
}

// GetPoint 获得当前点往下相交的navmesh的交点以及区域信息
func (m *Mesh) GetPointAndNode(pos linmath.Vector3) (linmath.Vector3, *Node, error) {

	ns, err := m.getNodeIDByPos(pos.X, pos.Z)
	if err != nil {
		return linmath.Vector3_Zero(), nil, err
	}

	for _, n := range ns {
		v, _, err := n.getInteractPoint(pos)
		if err == nil {
			return v, n, nil
		}
	}

	return linmath.Vector3_Zero(), nil, errors.New("no interact point")
}

func (m *Mesh) getNodeIDByPos(x, z float32) ([]*Node, error) {

	ix := int32(x)
	iz := int32(z)

	xm, ok := m.NodePosMap[ix]
	if !ok {
		return nil, errors.New("no node")
	}

	nl, ok := xm[iz]
	if !ok {
		return nil, errors.New("no node")
	}

	ns := make([]*Node, 0, 3)
	for _, n := range nl {
		if m.Nodes[n].isInclude(x, z) {
			ns = append(ns, m.Nodes[n])
		}
	}

	if len(ns) == 0 {
		return nil, errors.New("get node id data error")
	}

	//得到的node按照y轴坐标由上到下排序
	if len(ns) > 1 {
		m.sortNodeByHeight(ns)
	}

	return ns, nil
}

func (m *Mesh) sortNodeByHeight(ns []*Node) {
	for i := 1; i < len(ns); i++ {
		j := i
		for j > 0 && ns[j].center.Y > ns[j-1].center.Y {

			o := ns[j]
			ns[j] = ns[j-1]
			ns[j-1] = o

			j--
		}
	}
}

func (m *Mesh) getNeighbour(n *Node, b *Border) *Node {

	if b.NodeA == n.id && b.NodeB != -1 && m.Nodes[b.NodeB].area < 20 {
		return m.Nodes[b.NodeB]
	}

	if b.NodeB == n.id && b.NodeA != -1 && m.Nodes[b.NodeA].area <= 20 {
		return m.Nodes[b.NodeA]
	}

	return nil
}

func (m *Mesh) getCorner(cn *Node, nn *Node) (linmath.Vector3, linmath.Vector3) {

	var nodes int64
	if cn.id < nn.id {
		nodes = int64(cn.id)<<32 | int64(nn.id)
	} else {
		nodes = int64(nn.id)<<32 | int64(cn.id)
	}

	border, ok := m.BordersByNodes[nodes]
	if !ok {
		log.Error("get corner vert failed ")
	}

	if cn.id == border.NodeA && nn.id == border.NodeB {
		return m.Verts[border.VertA], m.Verts[border.VertB]
	} else if cn.id == border.NodeB && nn.id == border.NodeA {
		return m.Verts[border.VertB], m.Verts[border.VertA]
	}

	log.Error("get corner data error ")

	return linmath.Vector3_Zero(), linmath.Vector3_Zero()
}
