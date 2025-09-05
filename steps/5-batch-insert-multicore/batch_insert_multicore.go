package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Cria a estrutura do objeto fatura
type Fatura struct {
	emitente      string
	documento     string
	contrato      string
	categoria     string
	qtdNota       int
	fatura        int
	valor         float64
	dataCompra    string
	dataPagamento string
}

// Função responsável por gerar os dados para o canal, executada em goroutine
func dataProducer(numRegistros, chunkSize int, dataChan chan<- []Fatura) {
	// Fecha o canal quando a função terminar
	// Sinaliza aos workers que não precisam aguardar mais dados
	defer close(dataChan)

	// Separa as fatias de dados de acordo com o número de registros e tamanho da fatia (chunk)
	for i := 0; i < numRegistros; i += chunkSize {
		end := i + chunkSize
		if end > numRegistros {
			end = numRegistros
		}

		// Cria os dados de acordo com a fatia e envia para a fila do canal
		chunk := make([]Fatura, end-i)
		for j := 0; j < len(chunk); j++ {
			chunk[j] = Fatura{
				emitente:      "Blablabla Soluções Digitais",
				documento:     "401678cc1d9a716baebbac87452f62686572f729677c24921b3f08925d468e4c",
				contrato:      "17/2024",
				categoria:     "Prestação de serviços",
				qtdNota:       4,
				fatura:        34239 + i + j,
				valor:         20202953.58,
				dataCompra:    "2024-11-26",
				dataPagamento: "2025-01-30",
			}
		}
		dataChan <- chunk
	}
}

// Função que cada goroutine trabalhadora irá executar
func workerBatch(ctx context.Context, pool *pgxpool.Pool, dataChan <-chan []Fatura, wg *sync.WaitGroup) {
	// Sinaliza ao WaitGroup quando a função terminou de executar
	defer wg.Done()

	// Observa o canal dataChan e executa quando um novo dado é gerado
	for faturas := range dataChan {
		// Busca uma conexão do pool de conexões disponível no context
		conn, err := pool.Acquire(ctx)
		if err != nil {
			log.Printf("Erro ao adquirir conexão do pool: %v", err)
			continue
		}

		// Percorre as faturas da fatia e cria o batch statement para inserir
		batch := &pgx.Batch{}
		for _, fatura := range faturas {
			batch.Queue(
				`INSERT INTO faturas (emitente, documento, contrato, categoria, qtd_nota, fatura, valor, data_compra, data_pagamento) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				fatura.emitente, fatura.documento, fatura.contrato, fatura.categoria, fatura.qtdNota, fatura.fatura, fatura.valor, fatura.dataCompra, fatura.dataPagamento,
			)
		}

		// Envia o batch para o banco de dados
		br := conn.SendBatch(ctx, batch)
		br.Close()
		conn.Release()
	}
}

func main() {
	const connStr = "postgresql://postgres:123456@localhost:5432/postgres?sslmode=disable"
	const totalRegistros = 50000
	const chunkSize = 5000

	// Vou manter 4 workers no teste local
	// Na solução final, vou controlar isso na composição do docker
	// numWorkers := runtime.NumCPU()
	numWorkers := 4
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Cria o canal de dados que será usado pelo producer e workers
	dataChan := make(chan []Fatura)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Erro ao analisar a string de conexão: %v", err)
	}

	// Cria pool de conexões de acordo com o número de workers
	config.MaxConns = int32(numWorkers)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erro ao criar o pool de conexões: %v", err)
	}
	defer pool.Close()

	log.Println("Pool de conexões estabelecido. Criando tabela...")
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS faturas (
        id SERIAL PRIMARY KEY,
		emitente TEXT,
		documento TEXT,
		contrato TEXT,
		categoria TEXT,
		qtd_nota INT,
		fatura INT,
		valor DECIMAL(18, 2),
		data_compra DATE,
		data_pagamento DATE
	);`
	if _, err = pool.Exec(context.Background(), createTableSQL); err != nil {
		log.Fatalf("Erro ao criar a tabela: %v\n", err)
	}
	log.Println("Tabela 'faturas' criada ou já existente.")

	log.Printf("Iniciando a inserção de %d registros em %d goroutines...", totalRegistros, numWorkers)

	// Cria um controle de execução para as goroutines
	var wg sync.WaitGroup
	startTime := time.Now()

	// Cria uma go routine para o producer que alimenta o canal
	go dataProducer(totalRegistros, chunkSize, dataChan)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		// Cria uma go routine para o workder
		go workerBatch(context.Background(), pool, dataChan, &wg)
	}

	// Garante que a main espere por todos os trabalhadores terminarem o serviço antes de encerrar
	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("Inserção em fluxo (Batch) concluída com sucesso em %v", duration)
}
