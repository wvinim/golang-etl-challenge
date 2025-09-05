package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	BatchInsert()
}

func BatchInsert() {
	const connStr = "postgresql://postgres:123456@localhost:5432/postgres?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v\n", err)
	}
	// Agenda o fechamento da conexão ao final da execução da func
	defer conn.Close(context.Background())

	// Criação da tabela com exemplo simples de índices
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
	);
	CREATE INDEX IF NOT EXISTS idx_faturas_emitente ON faturas (emitente);
	CREATE INDEX IF NOT EXISTS idx_faturas_data_compra ON faturas (data_compra);
	CREATE INDEX IF NOT EXISTS idx_faturas_emitente_data_compra ON faturas (emitente, data_compra);`

	_, err = conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar a tabela: %v\n", err)
	}

	dados := struct {
		emitente      string
		documento     string
		contrato      string
		categoria     string
		qtdNota       int
		fatura        int
		valor         float64
		dataCompra    string
		dataPagamento string
	}{
		emitente:      "Blablabla Soluções Digitais",
		documento:     "401678cc1d9a716baebbac87452f62686572f729677c24921b3f08925d468e4c",
		contrato:      "17/2024",
		categoria:     "Prestação de serviços",
		qtdNota:       4,
		fatura:        34239,
		valor:         20202953.58,
		dataCompra:    "2024-11-26",
		dataPagamento: "2025-01-30",
	}

	// Insere 50000 linhas em batch
	const numRegistros = 50000
	batch := &pgx.Batch{}
	for i := 0; i < numRegistros; i++ {
		batch.Queue(
			`INSERT INTO faturas (emitente, documento, contrato, categoria, qtd_nota, fatura, valor, data_compra, data_pagamento) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			dados.emitente,
			dados.documento,
			dados.contrato,
			dados.categoria,
			dados.qtdNota,
			dados.fatura,
			dados.valor,
			dados.dataCompra,
			dados.dataPagamento,
		)
	}

	br := conn.SendBatch(context.Background(), batch)
	// Agenda o fechamento da operação em batch ao final da execução da func
	defer br.Close()
}
