package main

import (
	"protoMsg"
	"sync"
	"time"
	"zeus/linmath"

	"zeus/iserver"

	log "github.com/cihub/seelog"
)

type Space struct {
	spaceID_            uint64
	stype               uint32 // 1:大地图  2:其它
	cells               *sync.Map
	checkTicker         *time.Ticker
	tree                *bspTree     //对应树
	rect                linmath.Rect //需要边界么？先保留吧
	changingCellNode    *sync.Map    //需要移动边界的cell node
	setNeedMoveInterval *time.Ticker //需要多长时间添加新的移动结点
	//moveBoudaryTicker   *time.Ticker //cell边界移动、拆分定时器,优先移动，再拆分
}

func (sp *Space) Init(xmin float64, xmax float64, ymin float64, ymax float64) {
	log.Debug("space init, space id = ", sp.getID())
	sp.stype = 1
	sp.cells = &sync.Map{}
	sp.changingCellNode = &sync.Map{}
	sp.rect.Xmin = xmin
	sp.rect.Xmax = xmax
	sp.rect.Ymin = ymin
	sp.rect.Ymax = ymax

	sp.checkTicker = time.NewTicker(10 * time.Second)
	sp.setNeedMoveInterval = time.NewTicker(30 * time.Second)
	//sp.moveBoudaryTicker = time.NewTicker(30 * time.Second)

}

func (sp *Space) borderMove() {

	//移动条件：
	//1:遍历所有互为兄弟的叶子结点，兄弟俩必须都为叶子,changingCell里也不是所有的node都是叶子 ，因为存在多个node指向同一个cell
	//2:left cell和right cell都没有满，并且移动后left cell的大小不超过right cell的大小(也可以用左右cell的负载来做判定)

	cellmap := make(map[uint64]*protoMsg.RectInfo)
	sp.changingCellNode.Range(
		func(leftN, rightN interface{}) bool {
			left := leftN.(*bspNode)
			right := rightN.(*bspNode)
			if left.cell.getRectSize() < right.cell.getRectSize() {
				if left.cell.isHoriSplit() && right.cell.isHoriSplit() {
					left.cell.modifyBorder(true, true)
					right.cell.modifyBorder(false, true)

				} else if !(left.cell.isHoriSplit()) && !(right.cell.isHoriSplit()) {
					left.cell.modifyBorder(true, false)
					right.cell.modifyBorder(false, false)

				}

				leftrectinfo := left.cell.getProtoMsgRect()
				cellmap[left.cell.getID()] = leftrectinfo
				sp.cellChangeNotify(left.cell.getSrvID(), left.cell.getID(), leftrectinfo)

				rightrectinfo := right.cell.getProtoMsgRect()
				sp.cellChangeNotify(right.cell.getSrvID(), right.cell.getID(), rightrectinfo)
			}
			return true
		})

}

func (sp *Space) doloop() {
	select {
	case <-sp.checkTicker.C:
		//优先移动，然后拆分
		//暂时注释，方便测试
		sp.borderMove()
		sp.cells.Range(
			func(k, cellV interface{}) bool {
				cell := cellV.(*Cell)
				cell.MainLoop()
				return true
			})

	//设置添加新结点的时间间隔
	//case <-sp.setNeedMoveInterval.C:
	//	sp.setAllChangingCell()

	default:

	}
}

func (sp *Space) setID(cellID uint64) {

	sp.spaceID_ = cellID
}

func (sp *Space) getID() uint64 {
	return sp.spaceID_
}

func (sp *Space) getRect() *linmath.Rect {
	return &sp.rect
}

func (sp *Space) newCell(srvID uint64) *Cell {

	id := GetSrvInst().FetchTempID()
	cell := &Cell{
		cellID:  id,
		inSpace: sp,
		load:    0,
	}

	cell.setSrvID(srvID)
	sp.cells.Store(id, cell)

	return cell
}

//设置cell属于哪个场景服务器
/*func (sp *Space) setCellOwner(cell *Cell, srvID uint64) {
	if cell != nil {
		cell.setSrvID(srvID)
	}
}*/

func (sp *Space) newCellByLoad(cells map[uint64]uint32) {

	//新建cell
	for id, _ := range cells {
		cell := &Cell{
			cellID: id,
		}
		sp.cells.Store(id, cell)

	}
}

func (sp *Space) updateCellData(cellload map[uint64]uint32) {
	for id, load := range cellload {
		if cell, ok := sp.cells.Load(id); ok {
			cell.(*Cell).updateLoad(load)
		}
	}
}

func (sp *Space) getTree() *bspTree {
	return sp.tree
}

func (sp *Space) getCellByPos(pos *protoMsg.Vector3) *Cell {
	var cell *Cell = nil
	sp.cells.Range(
		func(k, cellV interface{}) bool {
			cell = cellV.(*Cell)
			if (cell.getRect().Xmin < float64(pos.X)) && (float64(pos.X) < cell.getRect().Xmax) {
				return false
			}
			return true
		})
	return cell
}

//添加需要改变边界的兄弟cell
func (sp *Space) addChangingNode(leftnode *bspNode, rightnode *bspNode) {
	if leftnode != nil {
		sp.changingCellNode.Store(leftnode, rightnode)
	}
}

//拆分过程中导致原来是叶子结点的，现在变成非叶子结点，所以需要删除
func (sp *Space) deleteChangingNode(leftnode *bspNode) {
	if leftnode != nil {
		sp.changingCellNode.Delete(leftnode)
	}
}

//设置需要改变边界的所有cell
func (sp *Space) setAllChangingCell() {
	sp.tree.Traversal()
}

//cell change notify
//cell边界移动信息同步给cellapp
func (sp *Space) cellChangeNotify(srvID uint64, cellid uint64, rectinfo *protoMsg.RectInfo) {

	msgret := &protoMsg.CellBorderChangeNotify{
		SpaceID:  sp.getID(),
		CellID:   cellid,
		Rectinfo: rectinfo,
	}

	log.Debug("send cell borderchange msg to cellapp .... = CellID ", msgret.CellID, " srvID= ", srvID)
	if err := iserver.GetSrvInst().PostMsgToCell(srvID, 0, msgret); err != nil {
		log.Error(err)
	}
}

//cell change notify
//同步拆分后的cell给cellapp，cellapp上会新增左儿子结点cell
func (sp *Space) allocNewCellNotify(cell *Cell, rectinfo *protoMsg.RectInfo) {

	msgret := &protoMsg.CreateCellNotify{
		SpaceID: sp.getID(),
		//后面需要space边界时再赋值
		Srect:   nil,
		CellID:  cell.getID(),
		Crect:   rectinfo,
		MapName: "1",
	}

	//获取将要分配新cell的CellApp的srvID
	cellapp := GetSrvInst().getFreeCellApp()
	if cellapp != nil {
		log.Debug("allocNewCellNotify to cellapp .... = spaceID =  ", sp.getID(), " cellID= ", cell.getID(), " srvID= ", cellapp.getID())
		cell.setSrvID(cellapp.getID())
		cellapp.addCell(cell)
		if err := iserver.GetSrvInst().PostMsgToCell(cellapp.getID(), 0, msgret); err != nil {
			log.Error(err)
		}
		cellapp.setValid(true)
	}
}

//cell delete notify
func (sp *Space) deleteCellNotify(cell *Cell) {

	msgret := &protoMsg.DeleteCellNotify{
		SpaceID: sp.getID(),
		CellID:  cell.getID(),
	}

	if err := iserver.GetSrvInst().PostMsgToCell(cell.getSrvID(), 0, msgret); err != nil {
		log.Error(err)
	}

}
