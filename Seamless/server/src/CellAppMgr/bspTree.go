package main

import (
	log "github.com/cihub/seelog"
)

type bspNode struct {
	Left   *bspNode
	Right  *bspNode
	Parent *bspNode
	isLeaf bool  //是否是叶子结点
	isHori bool  //是否水平分割
	cell   *Cell //只有叶子节点才对应cell

}

func (bn *bspNode) getSrvID() uint64 {
	return bn.cell.getSrvID()
}

//是否是叶子结点
func (bn *bspNode) isLeafN() bool {
	return bn.isLeaf
}

//设置是否是叶子结点
func (bn *bspNode) setleafN(bIsLeaf bool) {
	bn.isLeaf = bIsLeaf
}

func (bn *bspNode) isLeftNode() bool {
	if bn.Parent.Left == bn {
		return true
	}
	return false
}

//自己的兄弟是否是叶子结点
func (bn *bspNode) isBrotherLeaf() bool {

	if bn.Parent.Left == bn {
		return bn.Parent.Right.isLeafN()
	} else {
		return bn.Parent.Left.isLeafN()
	}
}

func (bn *bspNode) getBrother() *bspNode {

	if bn.Parent.Left == bn {
		return bn.Parent.Right
	} else {
		return bn.Parent.Left
	}
}

func (bn *bspNode) checkSelfAndBrotherNodeLeaf() (*bspNode, *bspNode, bool) {
	//根结点
	if bn.Parent == nil {
		return nil, nil, false
	} else if (bn.Parent.Left == bn) && (bn.isLeafN()) && (bn.Parent.Right.isLeafN()) {
		return bn, bn.Parent.Right, true
	} else if (bn.Parent.Right == bn) && (bn.isLeafN()) && (bn.Parent.Left.isLeafN()) {
		return bn.Parent.Left, bn, true
	} else {
		if bn.Parent.Left == bn {
			return bn, bn.Parent.Right, false
		} else if bn.Parent.Right == bn {
			return bn.Parent.Left, bn, false
		}
	}

	return nil, nil, false
}

type bspTree struct {
	root   *bspNode //根节点
	layer  uint32   //初始树的深度
	space  *Space   //对应space
	isInit bool     //树是否初始化完成

}

func (bt *bspTree) newBspNode(cell *Cell) *bspNode {

	node := &bspNode{
		cell:   cell,
		Left:   nil,
		Right:  nil,
		Parent: nil,
		isLeaf: true,
		isHori: false,
	}
	return node
}

func (bt *bspTree) Init(xmin float64, xmax float64, ymin float64, ymax float64, space *Space, cell *Cell) {

	//初始创建一个结点的树
	if bt.root == nil {
		//bt.root = bt.createLeafNode()
		bt.root = bt.newBspNode(cell)
		bt.space = space
		bt.root.cell = cell
		cell.setBspNode(bt.root)
		bt.isInit = true

	}

}

func (bt *bspTree) merge(node *bspNode) {

	//只有叶子结点才会合并

	log.Debug(" cellid= ", node.cell.getID(), " merge....")

	//改变parent cell大小
	if node.Parent.cell.isHoriSplit() {
		node.Parent.cell.setRectYmin(node.Parent.Left.cell.getRect().Ymin)
	} else {
		node.Parent.cell.setRectXmin(node.Parent.Left.cell.getRect().Xmin)
	}

	//通知父亲结点 大小 变化
	bt.space.cellChangeNotify(node.Parent.cell.getSrvID(), node.Parent.cell.getID(), node.Parent.cell.getProtoMsgRect())

	bt.space.deleteCellNotify(node.cell)
	//删除自己
	bt.deleteNode(node)

	//删除兄弟结点
	bt.deleteNode(node.getBrother())

}

//删除结点，这里删除的肯定是叶子结点
func (bt *bspTree) deleteNode(node *bspNode) {
	node.Parent = nil
}

func (bt *bspTree) split(node *bspNode, xmin float64, xmax float64, ymin float64, ymax float64) *bspNode {

	//百分比
	//产生新的cell,先从左边十分之一的地方切一刀
	var data float64 = 0.2
	var isHori = false

	// 找出分割线，并且确定是水平分割还是垂直分割, 通过长和宽谁更长，来决定是水平切还是垂直切，目的是保证切的均匀
	if node.cell.getRectX() < node.cell.getRectY() {
		isHori = true
	}

	return bt.genNode(node, xmin, xmax, ymin, ymax, isHori, data)

}

//生成新的结点,也就是这个节点的两个儿子
func (bt *bspTree) genNode(node *bspNode, xmin float64, xmax float64, ymin float64, ymax float64, isHori bool, data float64) *bspNode {

	//根据data判断是leftNode大还是rightnode大，谁大谁继承父亲，包括cellid都和父亲一样
	//规定初始创建的二叉树结点，左边小，右边大(无论是水平分割还是垂直分割)，也就是右子树结点和父亲结点指向相同的cell
	//拆分后，也就是附近结点和右儿子结点的cell面积一样，缩小了，分了一部分给左儿子结点

	//新建左边Node
	//获得合适的cellapp
	cellapp := GetSrvInst().getFreeCellApp()
	leftCell := bt.space.newCell(cellapp.getID())
	leftNode := bt.newBspNode(leftCell)
	leftNode.Parent = node

	//设置cell拆分保护期
	leftCell.splitProtectCount = 2
	leftCell.isHori = isHori

	//右边Node拷贝父亲
	rightNode := bt.newBspNode(node.cell)

	if isHori {

		log.Debug(" cellid= ", node.cell.getID(), " Horizontal split ...,  gen Node left node cellid = ", leftCell.getID(), " cell srvID = ", leftCell.getSrvID(), " xmin =  ", xmin, " xmax= ", xmax, " ymin=", ymin, " ymax= ", ymin+(data*(ymax-ymin)))
		log.Debug(" cellid= ", node.cell.getID(), " cell srvID = ", node.cell.getSrvID(), " Horizontal split ...,  gen Node right node cellid = ", node.cell.getID(), " cell srvID = ", node.cell.getSrvID(), " xmin =  ", xmin, " xmax= ", xmax, " ymin= ", ymin+(data*(ymax-ymin)), " ymax=", ymax)

		//从下面的 1/n开始拆分
		leftCell.Init(xmin, xmax, ymin, ymin+(data*(ymax-ymin)), leftNode)
		node.cell.Init(xmin, xmax, ymin+(data*(ymax-ymin)), ymax, rightNode)

	} else {
		log.Debug(" cellid= ", node.cell.getID(), " vertical split ...,  gen Node left node cellid = ", leftCell.getID(), " cell srvID = ", leftCell.getSrvID(), " xmin =  ", xmin, " xmax= ", xmin+(data*(xmax-xmin)), " ymin=", ymin, " ymax= ", ymax)
		log.Debug(" cellid= ", node.cell.getID(), " vertical split ..., gen Node right node cellid = ", node.cell.getID(), " cell srvID = ", node.cell.getSrvID(), " xmin =  ", xmin+(data*(xmax-xmin)), " xmax= ", xmax, " ymin= ", ymin, " ymax=", ymax)

		//从左边的 1/n开始拆分
		leftCell.Init(xmin, xmin+(data*(xmax-xmin)), ymin, ymax, leftNode)
		node.cell.Init(xmin+(data*(xmax-xmin)), xmax, ymin, ymax, rightNode)

	}
	leftNode.cell = leftCell
	node.Left = leftNode

	rightNode.cell = node.cell
	node.Right = rightNode
	//设置cell拆分保护期
	rightNode.cell.splitProtectCount = 2
	rightNode.cell.isHori = isHori
	rightNode.Parent = node

	//重置cell.node为右儿子
	node.cell.node = rightNode

	//设置父亲结点为非叶子结点
	node.setleafN(false)

	//拆分完设置需要改变边界的cell
	bt.space.setAllChangingCell()

	//先更新父亲节点所在CellApp的cell
	bt.space.cellChangeNotify(node.cell.getSrvID(), node.cell.getID(), node.cell.getProtoMsgRect())

	//再把左儿子cell分配给另外一个新的CellApp
	bt.space.allocNewCellNotify(leftCell, leftCell.getProtoMsgRect())

	return node
}

//前序遍历
func (bt *bspTree) Preorder(node *bspNode) {
	if node != nil {
		bt.Preorder(node.Left)
		bt.Preorder(node.Right)

		//找出所有的叶子结点
		//自己是叶子结点并且自己的兄弟也是叶子结点

		//log.Debug("checkSelfAndBrotherNodeLeaf left cell ")
		left, right, ok := node.checkSelfAndBrotherNodeLeaf()

		if ok {
			bt.space.addChangingNode(left, right)
		} else {
			bt.space.deleteChangingNode(left)
		}
	}
}

func (bt *bspTree) Traversal() {
	//前序遍历
	bt.Preorder(bt.root)
}
