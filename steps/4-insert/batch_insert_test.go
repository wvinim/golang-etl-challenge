package main

import (
	"testing"
)

func BenchmarkBatchInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BatchInsert()
	}
}
