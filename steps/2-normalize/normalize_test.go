package main

import (
	"testing"
)

var input = "Blablabla Soluções Digitais com vários ÁÉÍÓÚ Ç ç ãõ!?!? -"

func BenchmarkNormalizeHybrid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normalizeHybrid(input)
	}
}
