package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"golang-etl-challenge/consumer"
	"golang-etl-challenge/database"
	"golang-etl-challenge/models"
	"golang-etl-challenge/producer"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis do arquivo .env.
	// Se o arquivo não existir, ou houver um erro nas variáveis, o programa irá encerrar.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	connStr := os.Getenv("CONN_STR")
	if connStr == "" {
		log.Fatalf("Missing env variable CONN_STR")
	}
	filePatch := os.Getenv("FILE_PATH")
	if filePatch == "" {
		log.Fatalf("Missing env variable FILE_PATH")
	}
	chunkSize, _ := strconv.Atoi(os.Getenv("CHUNK_SIZE"))
	if chunkSize == 0 {
		log.Fatalf("Missing env variable CHUNK_SIZE")
	}

	// Usará o número de cores que estiver disponível na instância
	// Em docker, vou limitar em 4
	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Cria o canal de dados que será usado pelo producer e workers
	dataChan := make(chan []models.Fatura)

	// Cria pool de conexões de acordo com o número de workers
	pool, err := pgxpool.NewWithConfig(context.Background(), database.GetConfig(connStr, numWorkers))
	if err != nil {
		log.Fatalf("Erro ao criar o pool de conexões: %v", err)
	}
	defer pool.Close()

	// Cria a tabela no banco
	if err := database.CreateFaturasTable(context.Background(), pool); err != nil {
		log.Fatalf("Erro ao criar a tabela: %v", err)
	}

	// Verifica se existe o arquivo a ser processado
	if _, err := os.Stat(filePatch); os.IsNotExist(err) {
		log.Fatalf("Arquivo '%s' não encontrado. Crie o arquivo e tente novamente.", filePatch)
	}

	// Cria um controle de execução para as goroutines
	var wg sync.WaitGroup
	startTime := time.Now()

	// Cria uma go routine para o producer que alimenta o canal
	go producer.FileDataProducer(filePatch, chunkSize, dataChan)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		// Cria uma go routine para o worker
		go consumer.WorkerCopyFrom(context.Background(), pool, dataChan, &wg)
	}

	// Garante que a main espere por todos os trabalhadores terminarem o serviço antes de encerrar
	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("Inserção a partir do arquivo concluída com sucesso em %v", duration)
}
