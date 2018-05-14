package main

import (
	"container/list"
	"fmt"
)

const (
	StateNum = 10
)

// EntityStates 实体状态快照集合
type EntityStates struct {
	states *list.List

	lastState   IEntityState
	cachedState IEntityState
}

func NewEntityStates() *EntityStates {
	return &EntityStates{
		states: list.New(),
	}
}

func (es *EntityStates) addEntityState(state IEntityState) {
	state.SetDirty(true)
	if es.states.Len() >= StateNum && !es.states.Front().Next().Value.(IEntityState).IsDirty() {
		es.cachedState = es.states.Remove(es.states.Front()).(IEntityState)
	}
	es.states.PushBack(state)
	es.lastState = state
}

// GetLastState 获取当前的状态快照
func (es *EntityStates) GetLastState() IEntityState {
	if es.lastState == nil {
		panic("the states have at least one element")
	}

	return es.lastState
}

func (es *EntityStates) reflushDirtyState() {
	for p := es.states.Back(); p != nil; p = p.Prev() {
		state := p.Value.(IEntityState)
		if !state.IsDirty() {
			break
		}

		state.SetDirty(false)
	}

	num := es.states.Len() - StateNum
	if num > 0 {
		for i := 0; i < num; i++ {
			es.states.Remove(es.states.Front())
		}
	}
}

func (es *EntityStates) reflushDirtyStateAndGetDelta() ([]byte, bool) {
	f := es.getFirstNoDirtyElem()
	es.reflushDirtyState()
	b := es.states.Back()

	return f.Value.(IEntityState).Delta(b.Value.(IEntityState))
}

func (es *EntityStates) getFirstNoDirtyElem() *list.Element {

	for b := es.states.Back(); b != nil; b = b.Prev() {
		if !b.Value.(IEntityState).IsDirty() {
			return b
		}
	}

	panic("wrong ring buffer")
}

// GetHistoryState 获取历史快照信息
func (es *EntityStates) GetHistoryState(timeStamp uint32) IEntityState {

	b := es.states.Back()
	if b != nil {
		is := b.Value.(IEntityState)
		if is.GetTimeStamp() == timeStamp {
			return is
		}
	}

	for p := b.Prev(); p != nil; p = p.Prev() {
		v := p.Value
		if v != nil {
			is := v.(IEntityState)
			if is.GetTimeStamp() == timeStamp {
				return is
			} else if is.GetTimeStamp() < timeStamp {
				prevIS := is
				nextV := p.Next().Value
				if v == nil {
					return prevIS
				}

				nextIS := nextV.(IEntityState)
				return es.calcDiff(prevIS, nextIS, timeStamp)
			}
		}
	}

	return nil
}

func (es *EntityStates) calcDiff(prev, next IEntityState, timeStamp uint32) IEntityState {
	ns := prev.Clone()
	ns.SetTimeStamp(timeStamp)

	prevPos := prev.GetPos()
	prevRota := prev.GetRota()
	nextPos := next.GetPos()
	nextRota := next.GetRota()
	difPer := float32(timeStamp-prev.GetTimeStamp()) / float32(next.GetTimeStamp()-prev.GetTimeStamp())
	nsPos := prevPos.Add(nextPos.Sub(prevPos).Mul(difPer))
	nsRota := prevRota.Add(nextRota.Sub(prevRota).Mul(difPer))

	fmt.Println(prevPos, nextPos)
	fmt.Println(difPer)
	fmt.Println(prev.GetTimeStamp(), next.GetTimeStamp(), timeStamp)
	fmt.Println(nsPos)

	ns.SetPos(nsPos)
	ns.SetRota(nsRota)
	return ns
}
