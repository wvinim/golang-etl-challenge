package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Versão ultra rápida com substituição por tabela fixa (Considerando o conjunto de caracteres tradicionais do PT-BR)
var table [256]byte

func init() {
	// inicializa a tabela vazia
	for i := range table {
		table[i] = 0
	}

	// letras minúsculas ASCII → maiúsculas
	for c := 'a'; c <= 'z'; c++ {
		table[c] = byte(c - 32)
	}

	// letras maiúsculas
	for c := 'A'; c <= 'Z'; c++ {
		table[c] = byte(c)
	}

	// números
	for c := '0'; c <= '9'; c++ {
		table[c] = byte(c)
	}

	// espaço
	table[' '] = ' '
	table['-'] = ' '

	// acentos mais comuns PT-BR
	for _, c := range "áàãâäÁÀÃÂÄ" {
		table[c] = 'A'
	}
	for _, c := range "éèêëÉÈÊË" {
		table[c] = 'E'
	}
	for _, c := range "íìîïÍÌÎÏ" {
		table[c] = 'I'
	}
	for _, c := range "óòõôöÓÒÕÔÖ" {
		table[c] = 'O'
	}
	for _, c := range "úùûüÚÙÛÜ" {
		table[c] = 'U'
	}
	table['ç'], table['Ç'] = 'C', 'C'
}

func normalizeHybrid(s string) string {
	// O strings.Builder é usado para construir a string final, minimizando alocações de memória.
	var b strings.Builder
	// A função Grow é usada para pré-alocar a capacidade necessária, otimizando a performance.
	b.Grow(len(s))

	// Percorre os caracteres da string e busca eles na tabela fixa para substituição
	// Considerando que a maioria dos caracteres do campo serão sempre ASCII
	// A estratégia de decodificação híbrida se mostrou mais rápida 130 ns/op
	for i := 0; i < len(s); i++ {
		// Considera que o índice do caractere é ASCII (1 byte)
		charIndex := s[i]
		// Caso seja não ASCII (2 bytes ou mais)
		if charIndex > 128 {
			// Decodifica antes de fazer o de-para
			r, size := utf8.DecodeRuneInString(s[i:])
			// Incrementa corretamente o índice do for
			i += (size - 1)
			// Define o novo índice do caractere
			charIndex = byte(r)
		}
		// Realiza o de/para com a tabela fixa e salva na string
		if out := table[charIndex]; out != 0 {
			b.WriteByte(out)
		}
	}
	return b.String()
}

func main() {
	fmt.Println(normalizeHybrid("Blablabla Soluções Digitais com vários ÁÉÍÓÚ Ç ç ãõ!?!? -ç"))
}
