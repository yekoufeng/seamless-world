package main

import (
	//"zeus/linmath"
	"common"
	"zeus/linmath"

	"protoMsg"

	log "github.com/cihub/seelog"
)

type Cell struct {
	//cellID
	cellID   uint64
	cellinfo common.CellInfo
	//负载
	load uint32
	//所在space
	inSpace *Space
	//对应的node
	node *bspNode
	//拆分保护期，就是拆分后,多长时间不能再拆分,2个mainloop周期
	splitProtectCount uint32
	//记录这个cell是水平拆分还是垂直拆分的
	isHori bool
}

func (c *Cell) Init(xmin float64, xmax float64, ymin float64, ymax float64, node *bspNode) {

	c.cellinfo.Init(c.getID(), xmin, xmax, ymin, ymax)
	c.node = node
	c.splitProtectCount = 0
}

func (c *Cell) modifyBorder(isLeft bool, isHori bool) {
	//无论是水平切割还是垂直切割 都是left变大，right变小,先写死每个mainloop移动5
	log.Debug("cell modifyBorder cellid = ", c.getID(), " isLeft = ", isLeft, " isHori = ", isHori)
	if isHori {
		if isLeft {
			c.cellinfo.Init(c.getID(), c.getRect().Xmin, c.getRect().Xmax, c.getRect().Ymin, c.getRect().Ymax+5)
		} else {
			c.cellinfo.Init(c.getID(), c.getRect().Xmin, c.getRect().Xmax, c.getRect().Ymin+5, c.getRect().Ymax)
		}
	} else {
		if isLeft {
			c.cellinfo.Init(c.getID(), c.getRect().Xmin, c.getRect().Xmax+5, c.getRect().Ymin, c.getRect().Ymax)

		} else {
			c.cellinfo.Init(c.getID(), c.getRect().Xmin+5, c.getRect().Xmax, c.getRect().Ymin, c.getRect().Ymax)
		}
	}

}
func (c *Cell) getID() uint64 {
	return c.cellID
}

func (c *Cell) setSrvID(srvID uint64) {
	c.cellinfo.SetCellSrvID(srvID)
}
func (c *Cell) getSrvID() uint64 {
	return c.cellinfo.GetCellSrvID()
}

func (c *Cell) canSplit() bool {

	if c.splitProtectCount == 0 {
		return true
	} else {
		return false
	}
}

func (c *Cell) isHoriSplit() bool {
	//是水平拆分还是垂直拆分
	return c.isHori
}

//更新负载
func (c *Cell) updateLoad(load uint32) {
	//log.Debug("cell updateload cellid = ", c.getID(), " InServer= ", c.getSrvID(), " load = ", load)
	c.load = load
}

func (c *Cell) IsDestroyed() bool {
	return false
}

func (c *Cell) getTree() *bspTree {
	return c.inSpace.getTree()

}

func (c *Cell) getRect() *linmath.Rect {
	return c.cellinfo.GetRect()
}

func (c *Cell) getProtoMsgRect() *protoMsg.RectInfo {

	var rectinfo protoMsg.RectInfo
	rectinfo.Xmin = c.getRect().Xmin
	rectinfo.Xmax = c.getRect().Xmax
	rectinfo.Ymin = c.getRect().Ymin
	rectinfo.Ymax = c.getRect().Ymax
	return &rectinfo

}

func (c *Cell) getRectSize() float64 {
	rect := c.cellinfo.GetRect()
	return (rect.Xmax - rect.Xmin) * (rect.Ymax - rect.Ymin)
}

func (c *Cell) getRectX() float64 {
	rect := c.cellinfo.GetRect()
	return rect.Xmax - rect.Xmin
}

func (c *Cell) getRectY() float64 {
	rect := c.cellinfo.GetRect()
	return rect.Ymax - rect.Ymin
}

func (c *Cell) setRectXmin(xmin float64) {
	rect := c.cellinfo.GetRect()
	rect.Xmin = xmin
}

func (c *Cell) setRectYmin(ymin float64) {
	rect := c.cellinfo.GetRect()
	rect.Ymin = ymin
}

func (c *Cell) setRectXmax(xmax float64) {
	rect := c.cellinfo.GetRect()
	rect.Xmax = xmax
}

func (c *Cell) setRectYmax(ymax float64) {
	rect := c.cellinfo.GetRect()
	rect.Ymax = ymax
}

func (c *Cell) setBspNode(node *bspNode) {
	c.node = node
}

//处理拆分
func (c *Cell) dealSplit() {

	if c.load > common.SplitEntityNum {
		//log.Debug("cell MainLoop split..... cellid= ", c.getID(), " c.load= ", c.load, " splitProtectCount = ", c.splitProtectCount)

		if (c.getTree() != nil) && (c.node != nil) && (c.getTree().isInit) && (c.getRectSize() > 100) /*temp*/ && (c.canSplit()) {
			log.Debug("cell split......, c.cellid = ", c.getID(), " in server = ", c.getSrvID(), " c.load = ", c.load, " splitProtectCount = ", c.splitProtectCount, " rectX=", c.getRectX(), " rectY=", c.getRectY())
			c.getTree().split(c.node, c.getRect().Xmin, c.getRect().Xmax, c.getRect().Ymin, c.getRect().Ymax)
			c.load -= 20

		}
	}
	//递减拆分保护期
	if c.splitProtectCount > 0 {
		c.splitProtectCount -= 1
	}
}

//处理合并
func (c *Cell) dealMerge() {
	if (c.getTree() != nil) && (c.node != nil) && (c.getTree().isInit) && (c.getRectSize() < 100) && (c.load < common.MergeEntityNum) && (c.node.getBrother().cell.load < common.BeMergeEntityNum) {

		//必须是左儿子才能merge
		if c.node.Parent.Left == c.node {
			c.getTree().merge(c.node)
		}
	}
}

func (c *Cell) MainLoop() {

	c.dealSplit()
	c.dealMerge()
}
