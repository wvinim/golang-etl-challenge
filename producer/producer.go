package producer

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"

	"golang-etl-challenge/models"
	"golang-etl-challenge/normalizer"
)

// Lê os dados de um io.Reader, cria os lotes e os envia para o canal
// Esta função é agnóstica à fonte de dados
// Pode lidar com dados de arquivos, streams de rede, etc...
func FileDataProducer(reader io.Reader, chunkSize int, dataChan chan<- []models.Fatura) {
	// Fecha o canal quando a função terminar
	// Sinaliza aos workers que não precisam aguardar mais dados
	defer close(dataChan)

	// Define buffer máximo de 64KB
	// O melhor em diversos testes de performance (testei de 1 KB até 1 MB)
	scanner := bufio.NewScanner(reader)
	const maxCapacity = 64 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// Pulando o cabeçalho do arquivo (nomes das colunas)
	scanner.Scan()

	// Cria os dados de acordo com a fatia
	chunk := make([]models.Fatura, 0, chunkSize)
	linha := 2
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, "\t")

		if len(cols) != 9 {
			log.Printf("Linha com colunas fora do padrão, pulando: %d", linha)
			linha++
			continue
		}
		qtdNota, _ := strconv.Atoi(cols[4])
		fatura, _ := strconv.Atoi(cols[5])

		f := models.Fatura{
			Emitente:      normalizer.NormalizeCharactersHybrid(cols[0]),
			Documento:     normalizer.StringOrNil(cols[1]),
			Contrato:      normalizer.StringOrNil(cols[2]),
			Categoria:     normalizer.NormalizeCharactersHybrid(cols[3]),
			QtdNota:       qtdNota,
			Fatura:        fatura,
			Valor:         normalizer.NormalizeStringToFloat(cols[6]),
			DataCompra:    normalizer.NormalizeDateString(cols[7]),
			DataPagamento: normalizer.NormalizeDateString(cols[8]),
		}

		chunk = append(chunk, f)

		if len(chunk) >= chunkSize {
			dataChan <- chunk
			chunk = make([]models.Fatura, 0, chunkSize)
		}
		linha++
	}

	// Adiciona essa fatia ao canal de dados principal
	if len(chunk) > 0 {
		dataChan <- chunk
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erro durante leitura do arquivo: %v", err)
	}
}
