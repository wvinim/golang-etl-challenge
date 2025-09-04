package main

import (
	"testing"
)

func BenchmarkFileParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := FileParsing("../../resources/base_ficticia_dados_prova.txt")
		if err != nil {
			b.Fatalf("Erro: %v", err)
		}
	}
}
