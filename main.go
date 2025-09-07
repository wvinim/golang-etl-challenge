package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

const (
	// O tamanho padrão do lote.
	defaultChunkSize = 5000

	// A porta padrão do servidor
	defaultServerPort = ":8088"
)

var (
	// Variáveis compartilhadas entre as funcs
	pool       *pgxpool.Pool
	numWorkers int
	chunkSize  int
	serverPort string
)

// Define as variáveis de escopo geral
// Cria o pool de conexões e a tabela no banco
func init() {
	var err error

	environment := os.Getenv("ENVIRONMENT")
	// O docker importa as variáveis de ambiente do arquivo .env
	// Em ambiente local, é necessário usar o godotenv
	err = godotenv.Load(".env")
	if err != nil && environment != "docker" {
		log.Println("Erro ao carregar o arquivo .env, usando variáveis de ambiente do sistema.")
	}

	pgHost := os.Getenv("POSTGRES_HOST")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDb := os.Getenv("POSTGRES_DB")
	pgPort := os.Getenv("POSTGRES_PORT")

	if pgHost == "" || pgUser == "" || pgPassword == "" || pgDb == "" || pgPort == "" {
		log.Fatalf("Variáveis de ambiente de credenciais do Postgres ausentes.")
	}

	var serverPortEnv, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if serverPortEnv == 0 {
		serverPort = defaultServerPort
	} else {
		serverPort = fmt.Sprintf(":%d", serverPortEnv)
	}

	// Constrói a string de conexão.
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDb)

	// Define o chunkSize via env ou padrão
	chunkSize, _ = strconv.Atoi(os.Getenv("CHUNK_SIZE"))
	if chunkSize == 0 {
		chunkSize = defaultChunkSize
	}

	// Usará o número de cores que estiver disponível na instância
	// Em docker, vou limitar em 4
	numWorkers = runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Cria pool de conexões de acordo com o número de workers
	pool, err = pgxpool.NewWithConfig(context.Background(), database.GetConfig(connStr, numWorkers))
	if err != nil {
		log.Fatalf("Erro ao criar o pool de conexões: %v", err)
	}

	// Cria a tabela no banco
	if err := database.CreateFaturasTable(context.Background(), pool); err != nil {
		log.Fatalf("Erro ao criar a tabela: %v", err)
	}
}

func main() {
	// Agenda o encerramento da pool de conexões ao encerrar a func
	defer pool.Close()

	// Define a rota de upload e o handler
	http.HandleFunc("/upload", uploadHandler)

	// Cria o serviço HTTP, registra no log e encerra a aplicação em caso de erro
	log.Printf("Serviço iniciado na porta %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido.", http.StatusMethodNotAllowed)
		return
	}

	// Faz o parse para multipart/form-data para lidar com arquivos grandes
	// Caso o arquivo seja maior que 32MB, utiliza um arquivo temporário no disco para
	// não sobir tudo em memória
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Erro ao processar o formulário: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo do formulário: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Cria um controle de execução para as goroutines
	var wg sync.WaitGroup
	startTime := time.Now()

	// Cria o canal de dados que será usado pelo producer e workers
	dataChan := make(chan []models.Fatura)

	// Cria uma go routine gerenciável para o producer que alimenta o canal
	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.FileDataProducer(file, chunkSize, dataChan)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		// Cria uma go routine para cada worker (1 por core)
		go consumer.WorkerCopyFrom(context.Background(), pool, dataChan, &wg)
	}

	// Garante que a main espere por todos os workers terminarem o serviço antes de encerrar
	wg.Wait()

	duration := time.Since(startTime)

	responseString := fmt.Sprintf("Arquivo processado com sucesso em %v", duration)
	log.Printf("%s", responseString)

	// Retorna a request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseString))
}
