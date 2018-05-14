package db

import (
	"encoding/json"
	"fmt"

	log "github.com/cihub/seelog"
)

const (
	IdipMail = "IdipMailInfo"
)

type idipMailUtil struct {
}

type IdipMailData struct {
	MailID       uint32 `json:"MailId"`       // 邮件ID
	SendTime     uint32 `json:"SendTime"`     // 发送时间
	MailTitle    string `json:"MailTitle"`    // 邮件标题
	MailContent  string `json:"MailContent"`  // 邮件内容
	MinLevel     uint16 `json:"MinLevel"`     // 最小领取等级（默认0）
	MaxLevel     uint32 `json:"MaxLevel"`     // 最大领取等级（默认0）
	ItemOneID    uint32 `json:"ItemOneId"`    // 道具ID1
	ItemOneNum   uint32 `json:"ItemOneNum"`   // 道具数量1
	ItemTwoID    uint32 `json:"ItemTwoId"`    // 道具ID2
	ItemTwoNum   uint32 `json:"ItemTwoNum"`   // 道具数量2
	ItemThreeID  uint32 `json:"ItemThreeId"`  // 道具ID3
	ItemThreeNum uint32 `json:"ItemThreeNum"` // 道具数量3
	ItemFourID   uint32 `json:"ItemFourId"`   // 道具ID4
	ItemFourNum  uint32 `json:"ItemFourNum"`  // 道具数量4
	ItemFiveID   uint32 `json:"ItemFiveId"`   // 道具ID5
	ItemFiveNum  uint32 `json:"ItemFiveNum"`  // 道具数量5
	Hyperlink    string `json:"Hyperlink"`    // 超链接
	ButtonCon    string `json:"ButtonCon"`    // 按钮内容：(可以为空、为空时则不显示该超链接的按钮。不为空时则按钮显示输入的文字、如”点击查看“按钮)
}

type SliIdipMail []*IdipMailData

func (a SliIdipMail) Len() int           { return len(a) }
func (a SliIdipMail) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SliIdipMail) Less(i, j int) bool { return a[i].MailID < a[j].MailID }

func IdipMailUtil() *idipMailUtil {
	return &idipMailUtil{}
}

func (r *idipMailUtil) key() string {
	return fmt.Sprintf("%s", IdipMail)
}

func (r *idipMailUtil) GetMails(beginTime, endTime uint32) []*IdipMailData {
	res := []*IdipMailData{}
	info := hGetAll(r.key())
	log.Info("查询idip邮件数量", len(info))
	for _, v := range info {
		var d *IdipMailData
		if err := json.Unmarshal([]byte(v), &d); err != nil {
			log.Warn("Failed to Unmarshal ", err)
			return nil
		}

		if d.SendTime >= beginTime && d.SendTime <= endTime {
			res = append(res, d)
		}
	}
	return res
}

func (r *idipMailUtil) SaveMail(info *IdipMailData) {
	d, e := json.Marshal(info)
	if e != nil {
		log.Warn("marshal error ", e)
	}
	hSet(r.key(), info.MailID, string(d))
}

func (r *idipMailUtil) AddMail(info *IdipMailData) {
	r.SaveMail(info)
}
