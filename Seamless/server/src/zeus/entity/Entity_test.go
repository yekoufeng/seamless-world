package entity

import (
	"math/rand"
	"testing"
)

func BenchmarkRangeMapTypeA(b *testing.B) {
	props := make(map[int]*PropInfo)

	for i := 0; i < 10000; i++ {
		index := rand.Int()
		props[index] = &PropInfo{}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, p := range props {
			p = p
		}
	}
}

func BenchmarkRangeMapTypeB(b *testing.B) {
	props := make(map[int]*PropInfo)

	for i := 0; i < 10000; i++ {
		index := rand.Int()
		props[index] = &PropInfo{}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for index := range props {
			props[index] = nil
		}
	}
}
