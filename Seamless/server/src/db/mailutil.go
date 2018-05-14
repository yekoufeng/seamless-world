package db

import (
	"fmt"
	"protoMsg"
	"time"

	log "github.com/cihub/seelog"
)

const (
	mailPrefix             = "Mail"
	globalMailPrefix       = "globalMail"
	playerGlobalMailPrefix = "playerGlobalMailPrefix"
)

type mailUtil struct {
	uid uint64
}

func MailUtil(uid uint64) *mailUtil {
	return &mailUtil{
		uid: uid,
	}
}

func (r *mailUtil) key() string {
	return fmt.Sprintf("%s:%d", mailPrefix, r.uid)
}

func (r *mailUtil) GetMails() []*protoMsg.MailInfo {
	res := []*protoMsg.MailInfo{}
	info := hGetAll(r.key())
	for _, v := range info {
		var d protoMsg.MailInfo
		if err := d.Unmarshal([]byte(v)); err != nil {
			log.Warn("Failed to Unmarshal ", err)
			return nil
		}
		res = append(res, &d)
	}
	return res
}

func (r *mailUtil) GetMail(mailid uint64) *protoMsg.MailInfo {
	v := hGet(r.key(), mailid)
	var d protoMsg.MailInfo
	if err := d.Unmarshal([]byte(v)); err != nil {
		log.Warn("Failed to Unmarshal ", err)
		return nil
	}
	return &d
}

func (r *mailUtil) SaveMail(info *protoMsg.MailInfo) {
	d, e := info.Marshal()
	if e != nil {
		log.Warn("marshal error ", e)
	}
	hSet(r.key(), info.Mailid, string(d))
}

func (r *mailUtil) RemMail(mailid uint64) {
	hDEL(r.key(), mailid)
}

func (r *mailUtil) AddMail(info *protoMsg.MailInfo) {
	r.SaveMail(info)
}

func AddGlobalMail(info *protoMsg.MailInfo) {
	d, e := info.Marshal()
	if e != nil {
		log.Warn("marshal error ", e)
	}
	hSet(globalMailPrefix, info.Mailid, string(d))
}

func GetGlobalMails() []*protoMsg.MailInfo {
	res := []*protoMsg.MailInfo{}
	info := hGetAll(globalMailPrefix)
	for _, v := range info {
		var d protoMsg.MailInfo
		if err := d.Unmarshal([]byte(v)); err != nil {
			log.Warn("Failed to Unmarshal ", err)
			return nil
		}
		res = append(res, &d)
	}
	return res
}

type playerGlobalMailUtil struct {
	uid uint64
}

func PlayerGlobalMailUtil(uid uint64) *playerGlobalMailUtil {
	return &playerGlobalMailUtil{
		uid: uid,
	}
}

func (r *playerGlobalMailUtil) key() string {
	return fmt.Sprintf("%s:%d", playerGlobalMailPrefix, r.uid)
}

func (r *playerGlobalMailUtil) GetAll() map[string]string {
	info := hGetAll(r.key())
	return info
}

func (r *playerGlobalMailUtil) AddMail(mailid uint64) {
	hSet(r.key(), mailid, time.Now().Unix())
}
