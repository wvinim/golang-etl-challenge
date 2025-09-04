package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func FileParsing(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Define buffer máximo de 64KB
	// O melhor em diversos testes de performance (testei de 1 KB até 1 MB)
	const maxCapacity = 64 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	linha := 1
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, "\t")

		if len(cols) != 9 {
			// Exibindo possíveis linhas fora do padrão
			// A opção log.FatalLn() exibe o log e encerra a execução
			log.Println("Linha com mais colunas do que o esperado:", linha)
		}

		// fmt.Println("Emitente", cols[0])
		// fmt.Println("Documento", cols[1])
		// fmt.Println("Contrato", cols[2])
		// fmt.Println("Categoria", cols[3])
		// fmt.Println("qtdNota", cols[4])
		// fmt.Println("Fatura", cols[5])
		// fmt.Println("Valor", cols[6])
		// fmt.Println("data_compra", cols[7])
		// fmt.Println("data_pagamento", cols[8])
		// fmt.Println("-----------------------")
		linha++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erro durante leitura: %v", err)
	}
	return nil
}

func main() {
	if err := FileParsing("../../resources/base_ficticia_dados_prova.txt"); err != nil {
		log.Fatalf("Erro: %v", err)
	}
}
