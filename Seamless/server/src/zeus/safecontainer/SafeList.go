package safecontainer

import (
	"errors"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// SafeListNode 节点
type SafeListNode struct {
	next  unsafe.Pointer
	value interface{}
}

func newNode(data interface{}) unsafe.Pointer {
	return unsafe.Pointer(&SafeListNode{
		nil,
		data,
	})
}

// SafeList 安全链表
type SafeList struct {
	head unsafe.Pointer
	tail unsafe.Pointer

	C chan bool
}

// NewSafeList 新创建一个列表
func NewSafeList() *SafeList {

	node := unsafe.Pointer(newNode(nil))
	return &SafeList{
		node,
		node,
		make(chan bool, 1),
	}
}

// Put 放入
func (sl *SafeList) Put(data interface{}) {
	newNode := newNode(data)
	var tail unsafe.Pointer

	for {
		tail = atomic.LoadPointer(&sl.tail)
		next := atomic.LoadPointer(&(*SafeListNode)(tail).next)

		if next != nil {
			atomic.CompareAndSwapPointer(&sl.tail, tail, next)
		} else {
			if atomic.CompareAndSwapPointer(&(*SafeListNode)(sl.tail).next, nil, newNode) {
				break
			}
		}
		runtime.Gosched()
	}

	atomic.CompareAndSwapPointer(&sl.tail, tail, newNode)

	if len(sl.C) == 0 {
		sl.C <- true
	}
}

var errNoNode = errors.New("no node")

// Pop 拿出
func (sl *SafeList) Pop() (interface{}, error) {

	for {

		head := atomic.LoadPointer(&sl.head)
		tail := atomic.LoadPointer(&sl.tail)

		next := atomic.LoadPointer(&(*SafeListNode)(head).next)

		if head == tail {
			if next == nil {
				return nil, errNoNode
			}
			atomic.CompareAndSwapPointer(&sl.tail, tail, next)
		} else {
			if atomic.CompareAndSwapPointer(&sl.head, head, next) {
				return (*SafeListNode)(next).value, nil
			}
		}

		runtime.Gosched()
	}

}

// IsEmpty 是否为空
func (sl *SafeList) IsEmpty() bool {
	head := atomic.LoadPointer(&sl.head)
	tail := atomic.LoadPointer(&sl.tail)

	next := atomic.LoadPointer(&(*SafeListNode)(head).next)
	if head == tail {
		if next == nil {
			return true
		}
	}

	return false
}
