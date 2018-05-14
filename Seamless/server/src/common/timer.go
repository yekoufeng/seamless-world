package common

import (
	"container/heap"
	"time"

	"common/internals"
)

type Timer = internals.InnelTimer

type TimerMgr struct {
	impl *internals.InnelTimerMgr
}

func NewTimerMgr() *TimerMgr {
	return &TimerMgr{
		impl: &internals.InnelTimerMgr{},
	}
}

func (tm *TimerMgr) AddTimer(trigTime int64, f func()) *Timer {
	t := internals.NewInnelTimer(trigTime, f)
	heap.Push(tm.impl, t)
	return (*Timer)(t)
}

func (tm *TimerMgr) Tick() {
	if tm.impl.Len() == 0 {
		return
	}
	now := time.Now().UnixNano() / 1e6
	for tm.impl.Len() > 0 {
		first := tm.impl.First()
		if now < first.GetTrigTime() {
			break
		}
		first.Exec()
		heap.Pop(tm.impl)
	}
}

// 验证用的，实际用不到
func (tm *TimerMgr) CheckIndex() bool {
	return tm.impl.CheckIndex()
}

func (tm *TimerMgr) Remove(index int) {
	if index < 0 {
		return
	}
	heap.Remove(tm.impl, index)
}

func (tm *TimerMgr) Fix(index int) {
	if index < 0 {
		return
	}
	heap.Fix(tm.impl, index)
}
