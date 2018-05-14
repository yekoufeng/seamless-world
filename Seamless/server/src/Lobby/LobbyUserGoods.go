package main

import (
	"db"
	"excel"
	"protoMsg"
	"time"
	"zeus/iserver"

	log "github.com/cihub/seelog"
)

// StoreMgr 个人商店管理器
type StoreMgr struct {
	user *LobbyUser
}

func NewStoreMgr(user *LobbyUser) *StoreMgr {
	store := &StoreMgr{
		user: user,
	}
	return store
}

func (mgr *StoreMgr) initPropInfo() {

	data := db.PlayerGoodsUtil(mgr.user.GetDBID()).GetAllGoodsInfo()

	retMsg := &protoMsg.OwnGoodsInfo{
		List: make([]*protoMsg.OwnGoodsItem, 0),
	}

	for _, j := range data {
		item := &protoMsg.OwnGoodsItem{
			Id:    j.Id,
			State: j.State,
		}

		retMsg.List = append(retMsg.List, item)
	}

	// 初始化玩家拥有的所有购买商品
	//log.Debugf("初始化玩家拥有的所有购买商品 :%+v", retMsg)
	mgr.user.RPC(iserver.ServerTypeClient, "InitOwnGoodsInfo", retMsg)
}

// 商品类型
const (
	// 角色类型
	GoodsRoleType = 1
	// 伞包类型
	GoodsParachuteType = 2
	// 金币类型
	GoodsGoldCoin = 3
)

// 商品出售状态
const (
	// 出售中
	Onselling = 0
	// 不显示不出售
	NotShowNotSell = 1
	// 显示不出售
	ShowNotSell = 2
	// 免费
	GoodsFree = 3
)

// AddGoods 添加商品
func (mgr *StoreMgr) AddGoods(id uint32) bool {

	isOwn := db.PlayerGoodsUtil(mgr.user.GetDBID()).IsOwnGoods(id)
	if isOwn {
		//已经购买
		return false
	}

	info := &db.GoodsInfo{
		Id:    id,
		Time:  time.Now().Unix(),
		State: 0,
		Sum:   1,
	}

	// 添加至数据库
	if db.PlayerGoodsUtil(mgr.user.GetDBID()).AddGoodsInfo(info) == false {
		return false
	}

	// 通知客户端
	retMsg := &protoMsg.OwnGoodsItem{
		Id:    info.Id,
		State: info.State,
	}
	mgr.user.RPC(iserver.ServerTypeClient, "AddGoods", retMsg)

	return true
}

// UpdateGoodsState 更新物品状态
func (mgr *StoreMgr) updateGoodsState(id uint32, state uint32) {

	// 获取数据库信息
	info, err := db.PlayerGoodsUtil(mgr.user.GetDBID()).GetGoodsInfo(id)
	if err != nil || info == nil {
		log.Debug(err)
		return
	}

	// 更改数据库状态
	info.State = state
	result := db.PlayerGoodsUtil(mgr.user.GetDBID()).UpdateGoodsInfo(info)
	if !result {
		return
	}

	// 通知客户端
	retMsg := &protoMsg.OwnGoodsItem{
		Id:    info.Id,
		State: info.State,
	}
	mgr.user.RPC(iserver.ServerTypeClient, "UpdateGoods", retMsg)

}

// 邮件赠送商品
func (mgr *StoreMgr) MailGetGoods(id uint32, num uint32) {
	if num == 0 {
		return
	}

	// 商品信息
	goodsConfig, ok := excel.GetStore(uint64(id))
	if !ok {
		log.Error("邮件发送id错误 name:", mgr.user.GetName(), " id:", id)
		return
	}

	// 判断是否发放的是金币
	if goodsConfig.Type == GoodsGoldCoin {
		mgr.user.SetCoin(mgr.user.GetCoin() + uint64(num))
		return
	}

	// 是否拥有
	isOwn := db.PlayerGoodsUtil(mgr.user.GetDBID()).IsOwnGoods(id)
	if isOwn {

		// 获取数据库信息
		info, err := db.PlayerGoodsUtil(mgr.user.GetDBID()).GetGoodsInfo(id)
		if err != nil || info == nil {
			log.Debug(err)
			return
		}

		info.Sum += num
		db.PlayerGoodsUtil(mgr.user.GetDBID()).UpdateGoodsInfo(info)

	} else {

		info := &db.GoodsInfo{
			Id:    id,
			Time:  time.Now().Unix(),
			State: 0,
			Sum:   num,
		}

		result := db.PlayerGoodsUtil(mgr.user.GetDBID()).AddGoodsInfo(info)

		if result {
			// 通知客户端
			retMsg := &protoMsg.OwnGoodsItem{
				Id:    info.Id,
				State: info.State,
			}
			mgr.user.RPC(iserver.ServerTypeClient, "AddGoods", retMsg)
		}
	}

}
