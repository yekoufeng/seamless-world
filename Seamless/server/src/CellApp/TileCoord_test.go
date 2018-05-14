package main

// import (
// 	"fmt"
// 	"testing"
// 	"zeus/iserver"
// 	"zeus/linmath"
// )

// type TestCoordNode struct {
// 	ID        uint64
// 	Pos       linmath.Vector3
// 	IsNear    bool
// 	IsTrigger bool
// }

// func (n *TestCoordNode) GetID() uint64 {
// 	return n.ID
// }

// func (n *TestCoordNode) IsWatcher() bool {
// 	return true
// }

// func (n *TestCoordNode) SetPos(pos linmath.Vector3) {
// 	n.Pos = pos
// }
// func (n *TestCoordNode) GetPos() linmath.Vector3 {
// 	return n.Pos
// }

// func (n *TestCoordNode) IsNearAOILayer() bool {
// 	return n.IsNear
// }

// func (n *TestCoordNode) IsAOITrigger() bool {
// 	return n.IsTrigger
// }

// type TestAOICoordNode struct {
// 	TestCoordNode
// }

// func (n *TestAOICoordNode) OnEntityEnterAOI(o iserver.ICoordEntity) {
// 	fmt.Println(o.GetID(), "   enter to  ", n.GetID())
// }

// func (n *TestAOICoordNode) OnEntityLeaveAOI(o iserver.ICoordEntity) {
// 	fmt.Println(o.GetID(), " leave from ", n.GetID())
// }

// func TestCoord(t *testing.T) {

// 	n1 := &TestAOICoordNode{
// 		TestCoordNode{
// 			1,
// 			linmath.NewVector3(50, 0, 50),
// 			false,
// 			true,
// 		},
// 	}

// 	n2 := &TestAOICoordNode{
// 		TestCoordNode{
// 			2,
// 			linmath.NewVector3(150, 0, 50),
// 			false,
// 			true,
// 		},
// 	}

// 	n3 := &TestAOICoordNode{
// 		TestCoordNode{
// 			3,
// 			linmath.NewVector3(250, 0, 50),
// 			false,
// 			true,
// 		},
// 	}

// 	n4 := &TestAOICoordNode{
// 		TestCoordNode{
// 			4,
// 			linmath.NewVector3(50, 0, 150),
// 			false,
// 			true,
// 		},
// 	}

// 	n5 := &TestAOICoordNode{
// 		TestCoordNode{
// 			5,
// 			linmath.NewVector3(1500, 0, 1500),
// 			false,
// 			true,
// 		},
// 	}

// 	n6 := &TestAOICoordNode{
// 		TestCoordNode{
// 			6,
// 			linmath.NewVector3(250, 0, 150),
// 			false,
// 			true,
// 		},
// 	}

// 	n7 := &TestAOICoordNode{
// 		TestCoordNode{
// 			7,
// 			linmath.NewVector3(50, 0, 250),
// 			false,
// 			true,
// 		},
// 	}

// 	n8 := &TestAOICoordNode{
// 		TestCoordNode{
// 			8,
// 			linmath.NewVector3(150, 0, 250),
// 			false,
// 			true,
// 		},
// 	}

// 	n9 := &TestAOICoordNode{
// 		TestCoordNode{
// 			9,
// 			linmath.NewVector3(250, 0, 250),
// 			false,
// 			true,
// 		},
// 	}

// 	tile := NewTileCoord(10000, 10000)

// 	tile.UpdateCoord(n1)
// 	tile.UpdateCoord(n2)
// 	tile.UpdateCoord(n3)
// 	tile.UpdateCoord(n4)
// 	tile.UpdateCoord(n5)
// 	tile.UpdateCoord(n6)
// 	tile.UpdateCoord(n7)
// 	tile.UpdateCoord(n8)
// 	tile.UpdateCoord(n9)

// 	// n2.Pos = linmath.NewVector3(700, 0, 700)
// 	// tile.UpdateCoord(n2)
// 	fmt.Println("----------------------------------------")
// 	fmt.Println("----------------------------------------")
// 	n5.SetPos(linmath.NewVector3(150, 0, 150))
// 	tile.UpdateCoord(n5)
// }
