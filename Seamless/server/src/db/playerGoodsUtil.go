package db

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/cihub/seelog"
)

const (
	goodsPrefix = "OwnGoods"
)

// GoodsInfo 商品信息
type GoodsInfo struct {
	Id    uint32
	Time  int64
	State uint32
	Sum   uint32
}

// GoodsUtil
type GoodsUtil struct {
	uid uint64
}

func PlayerGoodsUtil(uid uint64) *GoodsUtil {
	return &GoodsUtil{
		uid: uid,
	}
}

func (p *GoodsUtil) key() string {
	return fmt.Sprintf("%s:%d", goodsPrefix, p.uid)
}

// GetAllGoodsInfo 获取所有商品信息
func (p *GoodsUtil) GetAllGoodsInfo() []*GoodsInfo {

	ret := make([]*GoodsInfo, 0)

	goodsList := hGetAll(p.key())

	for _, goods := range goodsList {

		var d *GoodsInfo
		if unErr := json.Unmarshal([]byte(goods), &d); unErr != nil {
			log.Error(unErr)
			continue
		}

		ret = append(ret, d)
	}

	return ret
}

// AddGoodsInfo 添加商品信息
func (p *GoodsUtil) AddGoodsInfo(info *GoodsInfo) bool {
	if info == nil {
		return false
	}

	if hExists(p.key(), info.Id) {
		return false
	}

	d, err := json.Marshal(info)
	if err != nil {
		log.Warn("AddGoodsInfo error = ", err)
		return false
	}

	hSet(p.key(), info.Id, string(d))
	return true
}

// UpdateGoodsInfo 更新商品信息
func (p *GoodsUtil) UpdateGoodsInfo(info *GoodsInfo) bool {

	if info == nil {
		return false
	}

	if !hExists(p.key(), info.Id) {
		return false
	}

	d, err := json.Marshal(info)
	if err != nil {
		log.Warn("UpdateGoodsInfo error = ", err)
		return false
	}

	hSet(p.key(), info.Id, string(d))
	return true
}

// IsBuyGoods 是否已拥有商品
func (p *GoodsUtil) IsOwnGoods(id uint32) bool {
	return hExists(p.key(), id)
}

// GetGoodsInfo 获取物品信息
func (p *GoodsUtil) GetGoodsInfo(id uint32) (*GoodsInfo, error) {

	if !hExists(p.key(), id) {
		return nil, errors.New("id is not exist")
	}

	var d *GoodsInfo
	v := hGet(p.key(), id)
	if err := json.Unmarshal([]byte(v), &d); err != nil {
		log.Warn("GetGoodsInfo Failed to Unmarshal ", err)
		return nil, errors.New("unmarshal error")
	}

	return d, nil

}
