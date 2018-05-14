package nav

import (
	"errors"
	"io/ioutil"
	"zeus/common"
	"zeus/linmath"
)

// NewMesh 根据一个路径生成一个mesh
func NewMesh(path string) *Mesh {

	mesh := &Mesh{
		Verts:          make([]linmath.Vector3, 0, 1000),
		Nodes:          make([]*Node, 0, 1000),
		BordersByNodes: make(map[int64]*Border),
		BordersByVerts: make(map[int64]*Border),
		NodePosMap:     make(map[int32]map[int32][]int32),
	}

	err := loadMesh(path, mesh)
	if err != nil {
		return nil
	}

	return mesh
}

func loadMesh(path string, m *Mesh) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	bs := common.NewByteStream(data)

	magicNum, _ := bs.ReadUInt32()
	if magicNum != 0xff335566 {
		return errors.New("wrong file header")
	}

	vertsNum, _ := bs.ReadUInt32()

	for i := 0; i < int(vertsNum); i++ {
		vert := linmath.Vector3_Zero()
		vert.X, _ = bs.ReadFloat32()
		vert.Y, _ = bs.ReadFloat32()
		vert.Z, _ = bs.ReadFloat32()

		m.Verts = append(m.Verts, vert)
	}

	bordersNum, _ := bs.ReadUInt32()

	for i := 0; i < int(bordersNum); i++ {
		border := &Border{}
		verts, _ := bs.ReadInt64()
		border.VertA, _ = bs.ReadInt32()
		border.VertB, _ = bs.ReadInt32()
		nodes, _ := bs.ReadInt64()
		border.NodeA, _ = bs.ReadInt32()
		border.NodeB, _ = bs.ReadInt32()
		border.Cost, _ = bs.ReadFloat32()

		m.BordersByNodes[nodes] = border
		m.BordersByVerts[verts] = border
	}

	nodesNum, _ := bs.ReadUInt32()

	for i := 0; i < int(nodesNum); i++ {
		node := &Node{
			id:      int32(i),
			borders: make([]*Border, 0, 5),
			poly:    make([]*linmath.Vector3, 0, 5),
		}

		vertsNum, _ := bs.ReadUInt32()

		for j := 0; j < int(vertsNum); j++ {
			vertID, _ := bs.ReadInt32()
			vert := &m.Verts[vertID]
			node.poly = append(node.poly, vert)
		}

		bordersNum, _ := bs.ReadUInt32()

		for j := 0; j < int(bordersNum); j++ {
			borderVert, _ := bs.ReadInt64()
			border := m.BordersByVerts[borderVert]
			node.borders = append(node.borders, border)
		}

		node.center.X, _ = bs.ReadFloat32()
		node.center.Y, _ = bs.ReadFloat32()
		node.center.Z, _ = bs.ReadFloat32()

		//node.normal.X, _ = bs.ReadFloat32()
		//node.normal.Y, _ = bs.ReadFloat32()
		//node.normal.Z, _ = bs.ReadFloat32()
		node.area, _ = bs.ReadInt32()

		m.Nodes = append(m.Nodes, node)
	}

	xNum, _ := bs.ReadUInt32()

	for i := 0; i < int(xNum); i++ {
		x, _ := bs.ReadInt32()

		xm := make(map[int32][]int32)
		m.NodePosMap[x] = xm

		yNum, _ := bs.ReadUInt32()

		for j := 0; j < int(yNum); j++ {
			y, _ := bs.ReadInt32()

			nodeNum, _ := bs.ReadUInt32()

			if nodeNum > 0 {

				nl := make([]int32, 0, nodeNum)

				for l := 0; l < int(nodeNum); l++ {
					n, _ := bs.ReadInt32()
					nl = append(nl, n)
				}
				xm[y] = nl
			}
		}
	}

	return nil
}
