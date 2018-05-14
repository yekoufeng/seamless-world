package online

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql" //mysql连接库
	"github.com/robfig/cron"
)

// Cnter 在线人数统计
type Cnter struct {
	srvID     uint64
	gameAppID string
	table     string
	zoneArea  int

	// 0:iOS, 1:Android
	onlineNum [2]int32

	mysqlDB  *sql.DB
	cronTask *cron.Cron
}

// NewCnter 新建在线人数统计
func NewCnter(user, pwd, addr, db, table, gameAppID string, zoneArea int, srvID uint64) (*Cnter, error) {
	cnt := &Cnter{}

	var err error

	cnt.mysqlDB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pwd, addr, db))
	if err != nil {
		return nil, err
	}
	cnt.gameAppID = gameAppID
	cnt.table = table
	cnt.srvID = srvID
	return cnt, nil
}

// Start 开始在线人数统计
func (cnt *Cnter) Start() error {
	cnt.cronTask = cron.New()
	if err := cnt.cronTask.AddFunc("0 */1 * * * ?", cnt.insertData); err != nil {
		return err
	}
	cnt.cronTask.Start()
	return nil
}

// Stop 停止在线人数统计
func (cnt *Cnter) Stop() error {
	cnt.cronTask.Stop()
	return cnt.mysqlDB.Close()
}

// ReportOnline 上报
func (cnt *Cnter) ReportOnline(plat int, num int32) {
	cnt.onlineNum[plat] += num
}

func (cnt *Cnter) insertData() {
	sqlStr := fmt.Sprintf("INSERT INTO %s(gameappid, timekey, gsid, zoneareaid, onlinecntios, onlinecntandroid) VALUES (?, ?, ?, ?, ?, ?)", cnt.table)
	_, err := cnt.mysqlDB.Exec(sqlStr, cnt.gameAppID, time.Now().Unix(), cnt.srvID, cnt.zoneArea, cnt.onlineNum[0], cnt.onlineNum[1])
	if err != nil {
		log.Error(err)
		return
	}
}
