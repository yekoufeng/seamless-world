package main

import (
	"common"
	"db"
	"excel"
	"protoMsg"
	"time"
	"zeus/iserver"

	log "github.com/cihub/seelog"
)

func (p *LobbyUserMsgProc) RPC_ReqGetMailList() {
	p.user.checkGlobalMail()
	mails := db.MailUtil(p.user.GetDBID()).GetMails()

	log.Debug(p.user.GetDBID(), "查看邮件列表", len(mails))
	if mails != nil && len(mails) != 0 {
		var msg protoMsg.RetMailList
		msg.Mails = mails
		p.user.RPC(iserver.ServerTypeClient, "RetMailList", &msg)
	}
}

func (p *LobbyUserMsgProc) RPC_ReqMailInfo(mailid uint64) {
	log.Info("查看邮件信息")
	mail := db.MailUtil(p.user.GetDBID()).GetMail(mailid)

	if mail != nil {
		mail.Haveread = true
		db.MailUtil(p.user.GetDBID()).SaveMail(mail)

		p.user.RPC(iserver.ServerTypeClient, "RetMailInfo", mail)
	}
}

func (p *LobbyUserMsgProc) RPC_DelMail(proto *protoMsg.DelMail) {
	log.Info("删除邮件")

	for _, v := range proto.Mailid {
		mail := db.MailUtil(p.user.GetDBID()).GetMail(v)
		if mail == nil {
			continue
		}
		if !mail.Haveget && len(mail.Objs) != 0 {
			continue
		}

		db.MailUtil(p.user.GetDBID()).RemMail(v)
	}
}

func (p *LobbyUserMsgProc) RPC_GetMailObj(mailid uint64) {
	log.Info("领取邮件附件", mailid)
	mail := db.MailUtil(p.user.GetDBID()).GetMail(mailid)

	if mail != nil {
		if mail.Haveget || len(mail.Objs) == 0 {
			p.user.RPC(iserver.ServerTypeClient, "GetMailObj", mailid, false)
			return
		}

		mail.Haveget = true
		db.MailUtil(p.user.GetDBID()).SaveMail(mail)

		for _, v := range mail.Objs {
			p.user.storeMgr.MailGetGoods(v.Id, v.Num)
		}

		p.user.RPC(iserver.ServerTypeClient, "GetMailObj", mailid, true)
	}
}

func checkMail(dbid uint64) {
	overdue := int64(common.GetTBSystemValue(common.System_MailOverdue))
	max := common.GetTBSystemValue(common.System_MailMax)
	now := time.Now().Unix()
	del := make([]uint64, 0)

	mails := db.MailUtil(dbid).GetMails()
	leftmail := make(map[uint64]*protoMsg.MailInfo)
	for _, v := range mails {
		if now >= int64(v.Gettime)+overdue {
			del = append(del, v.Mailid)
			continue
		}
		leftmail[v.Mailid] = v
	}

	if uint(len(leftmail)) > max {
		num := uint(len(leftmail)) - max
		var i uint
		for k, v := range leftmail {
			if v.Haveget || len(v.Objs) == 0 {
				delete(leftmail, k)
				del = append(del, k)
				i++
			}

			if i >= num {
				break
			}
		}
	}

	if uint(len(leftmail)) > max {
		num := uint(len(leftmail)) - max
		var i uint
		for k, _ := range leftmail {
			delete(leftmail, k)
			del = append(del, k)
			i++

			if i >= num {
				break
			}
		}
	}

	log.Info("删除邮件数量", len(del))
	for _, v := range del {
		db.MailUtil(dbid).RemMail(v)
	}
}

func (user *LobbyUser) MailNotify() {
	var havenew bool
	mails := db.MailUtil(user.GetDBID()).GetMails()
	for _, v := range mails {
		if !v.Haveread {
			havenew = true
		}
	}

	if havenew {
		user.RPC(iserver.ServerTypeClient, "AddNewMail")
	}
}

func (user *LobbyUser) checkGlobalMail() {
	overdue := int64(common.GetTBSystemValue(common.System_MailOverdue))
	now := time.Now().Unix()
	havesend := db.PlayerGlobalMailUtil(user.GetDBID()).GetAll()
	mails := db.GetGlobalMails()
	for _, v := range mails {
		mailid := common.Uint64ToString(v.Mailid)
		if _, ok := havesend[mailid]; !ok && now <= int64(v.Gettime)+overdue {
			db.PlayerGlobalMailUtil(user.GetDBID()).AddMail(v.Mailid)

			objs := make(map[uint32]uint32, 0)
			for _, obj := range v.Objs {
				objs[obj.Id] = obj.Num
			}
			sendObjMail(user.GetDBID(), v.Mailtype, v.Title, v.Text, v.Url, v.Button, objs)
		}
	}
}

func sendObjMail(dbid uint64, mailtype uint32, title string, text string, url string, button string, objs map[uint32]uint32) {
	mail := &protoMsg.MailInfo{}

	mailid := common.CreateNewMailID()
	if mailid == 0 {
		return
	}

	mail.Mailid = mailid
	mail.Mailtype = mailtype
	mail.Gettime = uint64(time.Now().Unix())
	mail.Haveread = false
	mail.Title = title
	mail.Text = text
	mail.Url = url
	mail.Button = button
	mail.Haveget = false

	for k, v := range objs {
		_, ok := excel.GetStore(uint64(k))
		if ok {
			obj := &protoMsg.MailObject{Id: k, Num: v}
			mail.Objs = append(mail.Objs, obj)
		}
	}

	db.MailUtil(dbid).AddMail(mail)
	log.Info("发送邮件", dbid, mailid, title)
}
