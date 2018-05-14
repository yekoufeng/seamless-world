package safecontainer

import (
	"fmt"
	"math/rand"
	"testing"
)

type msgFireInfo struct {
	name    string
	content interface{}
}

type T struct {
	index int
}

func Test_list(t *testing.T) {

	sl := NewSafeList()

	sl.Put(&msgFireInfo{
		name:    "test",
		content: &T{index: 1},
	})

	fmt.Println(sl.IsEmpty())

	sl.Pop()

	fmt.Println(sl.IsEmpty())
}

func BenchmarkList(b *testing.B) {
	sl := NewSafeList()

	b.ResetTimer()

	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			sl.Put(rand.Int())
			sl.Pop()
		}
	})

}

func BenchmarkChan(b *testing.B) {
	c := make(chan int, 1000)

	b.ResetTimer()

	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			c <- rand.Int()
			<-c
		}
	})
}
