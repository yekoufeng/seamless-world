package main

import (
	"common"
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"sync/atomic"
	"time"
	"zeus/zlog"

	"github.com/spf13/viper"

	log "github.com/cihub/seelog"
)

// BaseStaticInfo 总体统计信息
type BaseStaticInfo struct {
	Index         int
	RobotNumTotal int32
	RobotNumRun   int32
	RobotNumStop  int32
	Total         int64
	IntervalTotal int64

	startTime    int64
	intervalTime int64
	totalSec     int64
}

func (info *BaseStaticInfo) String() string {
	now := time.Now().Unix()
	intervalSec := now - info.intervalTime
	if intervalSec == 0 {
		intervalSec = 1
	}
	intervalTPS := info.IntervalTotal / intervalSec
	if info.IntervalTotal != 0 {
		info.totalSec += intervalSec
	}
	var tps int64
	if info.totalSec == 0 {
		tps = 0
	} else {
		tps = info.Total / info.totalSec
	}

	return fmt.Sprintf("date=%s|id=%d|type=0|robotNumTotal=%d|robotNumRun=%d|robotNumStop=%d|tps=%d|total=%d|totalSec=%d|intervalTPS=%d|intervalTotal=%d|intervalSec=%d\n",
		time.Now().Format("2006-01-02 15:04:05"), info.Index, info.RobotNumTotal,
		info.RobotNumRun, info.RobotNumStop, tps, info.Total, info.totalSec,
		intervalTPS, info.IntervalTotal, intervalSec,
	)
}

// TranStaticInfo 事务统计信息
type TranStaticInfo struct {
	Index         int
	Total         int64
	Err           int64
	TotalCostMS   int64
	MinMS         int64
	MaxMS         int64
	Per90MS       int64
	IntervalTotal int64
	IntervalErr   int64
	IntervalCost  int64
	IntervalMinMS int64
	IntervalMaxMS int64

	action       string
	startTime    int64
	intervalTime int64
	totalSec     int64
	totalSet     TranRetSet
}

// ResetInterval 重置统计间隔
func (info *TranStaticInfo) ResetInterval() {
	info.IntervalTotal = 0
	info.IntervalErr = 0
	info.IntervalCost = 0
	info.IntervalMinMS = math.MaxInt64
	info.IntervalMaxMS = 0
	info.intervalTime = time.Now().Unix()
}

func (info *TranStaticInfo) String() string {
	now := time.Now().Unix()

	intervalSec := now - info.intervalTime
	if intervalSec == 0 {
		intervalSec = 1
	}
	intervalTPS := info.IntervalTotal / intervalSec
	var intervalErrPer float32
	var intervalAvgMS int64
	if info.IntervalTotal != 0 {
		intervalErrPer = float32(info.IntervalErr / info.IntervalTotal)
		intervalAvgMS = info.IntervalCost / info.IntervalTotal
	}

	if info.IntervalTotal != 0 {
		info.totalSec += intervalSec
	}
	if info.totalSec == 0 {
		info.totalSec = 1
	}
	tps := info.Total / info.totalSec
	errPer := float32(info.Err / info.Total)
	avgMS := info.TotalCostMS / info.Total

	sort.Sort(info.totalSet)
	per90 := info.totalSet.Per90()
	info.Per90MS = per90.Cost

	if info.IntervalMinMS == math.MaxInt64 {
		info.IntervalMinMS = 0
	}
	if info.MinMS == math.MaxInt64 {
		info.MinMS = 0
	}
	return fmt.Sprintf("date=%s|id=%d|type=1|tps=%d|total=%d|err=%d|errPer=%0.2f%%|avgTimeMS=%d|minMillisec=%d|maxMillisec=%d|90PerMillisec=%d|intervalTPS=%d|intervalTotal=%d|intervalErr=%d|intervalErrPer=%0.2f%%|intervalAvgMS=%d|intervalMinMillisec=%d|intervalMaxMillisec=%d|action=%s\n",
		time.Now().Format("2006-01-02 15:04:05"), info.Index, tps, info.Total, info.Err,
		errPer, avgMS, info.MinMS, info.MaxMS, info.Per90MS, intervalTPS, info.IntervalTotal,
		info.IntervalErr, intervalErrPer, intervalAvgMS, info.IntervalMinMS, info.IntervalMaxMS,
		info.action,
	)
}

// TranRet 事务结果
type TranRet struct {
	Action string
	Err    error
	Cost   int64
}

// TranRetSet 事务结果集
type TranRetSet []*TranRet

func (set TranRetSet) Len() int {
	return len(set)
}

func (set TranRetSet) Less(i, j int) bool {
	return set[i].Cost < set[j].Cost
}

func (set TranRetSet) Swap(i, j int) {
	set[i], set[j] = set[j], set[i]
}

// Per90 返回90%的事务时间
func (set TranRetSet) Per90() *TranRet {
	index := len(set) * 9 / 10
	return set[index]
}

func main() {
	num := flag.Int("num", 5000, "客户端数量")
	action := flag.String("action", "reg", "注册场景")
	mapid := flag.Int("map", 1, "骷髅岛1, 致命丛林2, 默认2")
	flag.Parse()

	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Error("加载配置文件失败", err)
		return
	}

	zlog.Init("", "debug")
	defer log.Flush()

	ctx, ctxCancel := context.WithCancel(context.Background())

	id := viper.GetString("Config.ID")

	base := &BaseStaticInfo{}
	base.startTime = time.Now().Unix()
	base.intervalTime = base.startTime

	retChan := make(chan *TranRet, 1000)
	go HandleResult(retChan, base)
	common.InitMsg()
	tps := viper.GetInt("Config.TPS")
	tpsTicker := time.NewTicker(time.Duration(1000000/tps) * time.Microsecond)
	index := 0

	go func() {
		for {
			if index >= *num {
				break
			}
			index++

			select {
			case <-tpsTicker.C:
				atomic.AddInt32(&base.RobotNumTotal, 1)
				atomic.AddInt32(&base.RobotNumRun, 1)
				go func(i int) {
					defer func() {
						if err := recover(); err != nil {
							if viper.GetBool("Config.Recover") {
								panic(err)
							} else {
								fmt.Println(err)
							}
						}
					}()

					switch *action {
					case "reg":
						user := fmt.Sprintf("%s%d", id, time.Now().UnixNano())
						doRegister(user, "123", retChan)
					case "login":
						user := fmt.Sprintf("%s%d", id, i)
						doLoginAll(ctx, user, "123", retChan)
					case "relogin":
						user := fmt.Sprintf("%s%d", id, i)
						doLoginRepeat(ctx, user, "123", retChan)
					case "match":
						user := fmt.Sprintf("%s%d", id, i)
						doMatch(ctx, user, "123", retChan, false)
					case "team":
						user := fmt.Sprintf("%s%d", id, i)
						doMatch(ctx, user, "123", retChan, true)
					case "game":
						user := fmt.Sprintf("%s%d", id, i)
						// per := float32(index / *num)
						// if per < 0.01 {
						// 	doGame(ctx, user, "123", retChan, 4, uint32(*mapid))
						// } else if per < 0.03 {
						// 	doGame(ctx, user, "123", retChan, 3, uint32(*mapid))
						// } else if per < 0.06 {
						// 	doGame(ctx, user, "123", retChan, 2, uint32(*mapid))
						// } else {
						doGame(ctx, user, "123", retChan, 1, uint32(*mapid))
						// }
					case "regame":
						user := fmt.Sprintf("%s%d", id, i)
						per := float32(index / *num)
						if per < 0.01 {
							doGameRepeat(ctx, user, "123", retChan, 4, uint32(*mapid))
						} else if per < 0.03 {
							doGameRepeat(ctx, user, "123", retChan, 3, uint32(*mapid))
						} else if per < 0.06 {
							doGameRepeat(ctx, user, "123", retChan, 2, uint32(*mapid))
						} else {
							doGameRepeat(ctx, user, "123", retChan, 1, uint32(*mapid))
						}
					}
					atomic.AddInt32(&base.RobotNumStop, 1)
					atomic.AddInt32(&base.RobotNumRun, -1)
				}(index)
			}
		}

		fmt.Println("机器人启动完成")
	}()

	dur := viper.GetInt("Config.Time")
	timeOut := time.NewTimer(time.Duration(dur) * time.Second)
	select {
	case <-timeOut.C:
		fmt.Println("测试完成")
		ctxCancel()
		return
	}
}

// HandleResult 处理事务结果
func HandleResult(tranRet chan *TranRet, base *BaseStaticInfo) {
	targetFileName := fmt.Sprintf("stress-%s", time.Now().Format("2006-01-02-15-04-05"))
	f, err := os.Create(targetFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	trans := make(map[string]*TranStaticInfo)
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			base.Index++
			f.WriteString(base.String())
			base.intervalTime = time.Now().Unix()
			base.IntervalTotal = 0

			for _, tran := range trans {
				tran.Index++
				f.WriteString(tran.String())
				tran.ResetInterval()
			}
		case ret := <-tranRet:
			base.Total++
			base.IntervalTotal++

			tran, ok := trans[ret.Action]
			if !ok {
				tran = &TranStaticInfo{}
				tran.action = ret.Action
				now := time.Now().Unix()
				tran.startTime = now
				tran.intervalTime = now
				tran.totalSet = make([]*TranRet, 0, 1)
				tran.MinMS = math.MaxInt64
				tran.IntervalMinMS = math.MaxInt64
				trans[ret.Action] = tran
			}

			tran.Total++
			tran.IntervalTotal++
			tran.TotalCostMS += ret.Cost
			tran.IntervalCost += ret.Cost
			tran.totalSet = append(tran.totalSet, ret)
			if ret.Cost < tran.MinMS {
				tran.MinMS = ret.Cost
			}
			if ret.Cost > tran.MaxMS {
				tran.MaxMS = ret.Cost
			}
			if ret.Cost > tran.IntervalMaxMS {
				tran.IntervalMaxMS = ret.Cost
			}
			if ret.Cost < tran.IntervalMinMS {
				tran.IntervalMinMS = ret.Cost
			}
			if ret.Err != nil {
				tran.IntervalErr++
				tran.Err++
			}
		}
	}
}
