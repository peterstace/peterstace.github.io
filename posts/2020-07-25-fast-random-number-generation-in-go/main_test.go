package main

import (
	"math/rand"
	"testing"
)

var k = 100

func BenchmarkMathRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(k)
	}
}

func BenchmarkMathRandLocal(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rnd.Intn(k)
	}
}

func BenchmarkCreateLocalRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.New(rand.NewSource(0))
	}
}

var result int

func BenchmarkLCG(b *testing.B) {
	var cumulative int
	var state int32
	for i := 0; i < b.N; i++ {
		state = state*1664525 + 1013904223
		rndNum := int(state) % k
		cumulative += rndNum // prevent optimising everything out
	}
	result = cumulative
}

func BenchmarkFastLCG(b *testing.B) {
	var cumulative int
	var state int32
	for i := 0; i < b.N; i++ {
		state = state*1664525 + 1013904223
		rndNum := (int64(state) * int64(k)) >> 32
		cumulative += int(rndNum) // prevent optimising everything out
	}
	result = cumulative
}
