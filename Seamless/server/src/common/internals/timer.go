package internals

import "container/heap"

type InnelTimer struct {
	trigTime int64
	f        func()
	index    int
}

func NewInnelTimer(trigTime int64, f func()) *InnelTimer {
	return &InnelTimer{
		trigTime: trigTime,
		f:        f,
		index:    -1,
	}
}

func (t *InnelTimer) GetIndex() int {
	return t.index
}
func (t *InnelTimer) SetTrigTime(trigTime int64) {
	t.trigTime = trigTime
}

func (t *InnelTimer) GetTrigTime() int64 {
	return t.trigTime
}

func (t *InnelTimer) Exec() {
	t.f()
}

type InnelTimerMgr struct {
	timers []*InnelTimer
	heap   heap.Interface
}

func (tm *InnelTimerMgr) Len() int {
	return len(tm.timers)
}

func (tm *InnelTimerMgr) Less(i, j int) bool {
	return tm.timers[i].trigTime < tm.timers[j].trigTime
}

func (tm *InnelTimerMgr) Swap(i, j int) {
	tm.timers[i], tm.timers[j] = tm.timers[j], tm.timers[i]
	tm.timers[i].index = i
	tm.timers[j].index = j
}

func (tm *InnelTimerMgr) Push(x interface{}) {
	t := x.(*InnelTimer)
	t.index = len(tm.timers)
	tm.timers = append(tm.timers, t)
}

func (tm *InnelTimerMgr) Pop() interface{} {
	length := len(tm.timers)
	if length == 0 {
		return nil
	}

	v := tm.timers[length-1]
	tm.timers = tm.timers[0 : length-1]
	v.index = -1
	return v
}

func (tm *InnelTimerMgr) First() *InnelTimer {
	if len(tm.timers) == 0 {
		return nil
	}

	return tm.timers[0]
}

func (tm *InnelTimerMgr) CheckIndex() bool {
	for i := 0; i < len(tm.timers); i++ {
		if i != tm.timers[i].index {
			return false
		}
	}
	return true
}
