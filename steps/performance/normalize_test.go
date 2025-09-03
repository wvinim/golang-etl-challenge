package main

import (
	"testing"
)

var input = "Innovatech Soluções Digitais com vários ÁÉÍÓÚ Ç ç ãõ!?!?"

func BenchmarkNormalize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normalize(input)
	}
}

func BenchmarkNormalizeFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normalizeFast(input)
	}
}
