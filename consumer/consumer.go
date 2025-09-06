package consumer

import (
	"context"
	"log"
	"sync"

	"golang-etl-challenge/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Função que cada goroutine consumidora irá executar
func WorkerCopyFrom(ctx context.Context, pool *pgxpool.Pool, dataChan <-chan []models.Fatura, wg *sync.WaitGroup) {
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

		// Popula a fatia da fonte temporária com a interface compatível (CopyFrom)
		fs := models.NewFaturaSource(faturas)
		// Cria a lista de colunas da source
		columnNames := []string{"emitente", "documento", "contrato", "categoria", "qtd_nota", "fatura", "valor", "data_compra", "data_pagamento"}

		// Executa o CopyFrom para o banco de dados e libera a conexão
		if _, err = conn.CopyFrom(ctx, pgx.Identifier{"faturas"}, columnNames, fs); err != nil {
			log.Printf("Erro na operação CopyFrom: %v", err)
		}
		conn.Release()
	}
}
