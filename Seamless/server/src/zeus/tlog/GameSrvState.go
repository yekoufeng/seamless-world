package tlog

import "time"

// GameSvrState 服务器状态日志结构
type GameSvrState struct {
	DtEventTime string
	VGameIP     string
	IZoneAreaID uint32
}

// Name 结构体名字
func (state *GameSvrState) Name() string {
	return "GameSvrState"
}

// StateLogger 服务器状态日志记录
type StateLogger struct {
	state    *GameSvrState
	interval time.Duration
	stopC    chan bool
}

// NewStateLogger 创建新的服务器状态记录器
func NewStateLogger(addr string, zone uint32, interval time.Duration) *StateLogger {
	logger := &StateLogger{}
	logger.state = &GameSvrState{}
	logger.state.VGameIP = addr
	logger.state.IZoneAreaID = zone
	logger.interval = interval
	logger.stopC = make(chan bool, 1)
	return logger
}

// Start 启动记录器, 在协程中运行
func (logger *StateLogger) Start() {
	logger.state.DtEventTime = time.Now().Format("2006-01-02 15:04:05")
	Format(logger.state)
	go func() {
		ticker := time.NewTicker(logger.interval)
		defer ticker.Stop()
		for {
			select {
			case <-logger.stopC:
				return
			case <-ticker.C:
				logger.state.DtEventTime = time.Now().Format("2006-01-02 15:04:05")
				Format(logger.state)
			}
		}
	}()
}

// Stop 停止记录器
func (logger *StateLogger) Stop() {
	logger.stopC <- true
}
