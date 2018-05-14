package main

import "zeus/common"

//CreateRealEntity 创建real
func CreateRealEntity() *RealEntity {
	realEntity := &RealEntity{}
	realEntity.init()

	return realEntity
}

// Haunt 记录ghost所在的位置
type Haunt struct {
	serverid uint64
}

// RealEntity 真是的实体，与ghost相对
//作为real实体，需要记录一些额外数据
type RealEntity struct {
	hauntMap map[uint64]*Haunt
}

func (r *RealEntity) init() {
	r.hauntMap = make(map[uint64]*Haunt)
}

func (r *RealEntity) clearHaunts() {
	r.hauntMap = make(map[uint64]*Haunt)
}

func (r *RealEntity) packHaunt() (int, []byte) {
	size := len(r.hauntMap) * 16
	hauntBytes := make([]byte, size)
	bs := common.NewByteStream(hauntBytes)

	for cellID, haunt := range r.hauntMap {
		if err := bs.WriteUInt64(cellID); err != nil {
			//log.Error("Pack cellID failed ", err)
		}

		if err := bs.WriteUInt64(haunt.serverid); err != nil {
			//log.Error("Pack haunt failed ", err)
		}
	}

	return len(r.hauntMap), hauntBytes
}

func (r *RealEntity) umpackHaunt(num int, data []byte) {

	bs := common.NewByteStream(data)
	for i := 0; i < num; i++ {
		cellID, err := bs.ReadUInt64()
		if err != nil {
			//e.Error("read prop name fail ", err)
			return
		}

		haunt := &Haunt{}

		haunt.serverid, err = bs.ReadUInt64()
		if err != nil {
			//e.Error("read prop from stream failed ", name, err)
			return
		}

		r.hauntMap[cellID] = haunt
	}
}
