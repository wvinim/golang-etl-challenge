package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Cria a configuração do pool de conexões.
func GetConfig(connStr string, maxConns int) *pgxpool.Config {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Erro ao analisar a string de conexão: %v", err)
	}
	config.MaxConns = int32(maxConns)
	return config
}

// Cria a tabela 'faturas' no banco de dados.
func CreateFaturasTable(ctx context.Context, pool *pgxpool.Pool) error {
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

	_, err := pool.Exec(ctx, createTableSQL)
	if err != nil {
		return err
	}

	return nil
}
