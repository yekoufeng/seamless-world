package common

import (
	"protoMsg"
	"zeus/linmath"

	log "github.com/cihub/seelog"
)

//CellInfo
type CellInfo struct {
	rect   linmath.Rect
	srvID  uint64
	cellID uint64
}

func (info *CellInfo) Init(cellid uint64,  xmin float64, xmax float64, ymin float64, ymax float64) {
	info.cellID = cellid
	log.Debug("CellInfo init cellid  = ", info.cellID, " srvID = ", info.srvID, " xmin = ", xmin, " xmax = ", xmax, " ymin= ", ymin, " ymax=", ymax)
	info.rect.Init(xmin, xmax, ymin, ymax)
}

func (info *CellInfo) GetRect() *linmath.Rect {

	return &info.rect
}

//获得cell的width
func (info *CellInfo) GetCellLength() float64 {
	length := info.rect.Xmax - info.rect.Xmin
	if length < 0 {
		return -length
	} else {
		return length
	}
}

//获得cell的width
func (info *CellInfo) GetCellWidth() float64 {
	width := info.rect.Ymax - info.rect.Ymin
	if width < 0 {
		return -width
	} else {
		return width
	}
}

func (info *CellInfo) SetCellSrvID(srvID uint64) {
	info.srvID = srvID
}

//获得cell的serverID
func (info *CellInfo) GetCellSrvID() uint64 {
	return info.srvID
}

//获得cell的ID
func (info *CellInfo) GetCellID() uint64 {
	return info.cellID
}

func (info *CellInfo) FillRect(rect *protoMsg.RectInfo) {
	info.rect.Xmin = rect.Xmin
	info.rect.Xmax = rect.Xmax
	info.rect.Ymin = rect.Ymin
	info.rect.Ymax = rect.Ymax
}


