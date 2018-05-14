package timer

import (
	"fmt"
	"testing"
	"time"
)

type T struct {
}

func (t *T) TimerPrint() {
	fmt.Println("TimerPrint")
}

func (t *T) DelayPrint() {
	fmt.Println("DelayPrint")
}

func TestTimerObj(t *testing.T) {
	t1 := &T{}

	timer := NewTimer()

	timer.AddDelayCallByObj(t1, t1.DelayPrint, 1*time.Second)
	timer.RegTimerByObj(t1, t1.TimerPrint, 1*time.Second)

	time.Sleep(5 * time.Second)

	for {
		timer.Loop()

		time.Sleep(1 * time.Second)

		if err := timer.UnregTimerByObj(t1); err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkTimer(b *testing.B) {
	timer := NewTimer()

	for i := 0; i < 100000; i++ {
		timer.RegTimer(func1, 2*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timer.Loop()
	}
}

/*
func ExampleTimer() {
	timer := NewTimer()

	timer.RegTimer(func1, 2*time.Second)
	f2handle := timer.RegTimer(func2, 1*time.Second)

	looptime := 0
	for {
		if looptime >= 10 {
			return
		}

		if looptime == 5 {
			timer.SuspendTimer(f2handle)
		}

		if looptime == 8 {
			timer.ResumeTimer(f2handle)
		}

		timer.Loop()

		time.Sleep(1 * time.Second)
		looptime++
	}
}

func ExampleTimerOneTime() {
	timer := NewTimer()

	timer.AddDelayCall(func1, 3*time.Second)
	looptime := 0
	for {
		if looptime >= 10 {
			return
		}
		timer.Loop()
		time.Sleep(1 * time.Second)
		looptime++
	}

	// Output:
	// func1 Called
}
*/
func func1() {

}

func func2() {
	fmt.Println("func2 Called")
}

func func3() {
	fmt.Println("func3 Called")
}
