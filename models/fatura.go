package models

import (
	"github.com/jackc/pgx/v5"
)

// Representa a estrutura de um registro de fatura.
type Fatura struct {
	Emitente      *string
	Documento     *string
	Contrato      *string
	Categoria     *string
	QtdNota       int
	Fatura        int
	Valor         float64
	DataCompra    *string
	DataPagamento *string
}

// Representa a implementação de pgx.CopyFromSource para o CopyFrom.
type faturaSource struct {
	faturas []Fatura
	pos     int
}

// Representa o conjunto de dados já estruturados para o CopyFrom
func NewFaturaSource(faturas []Fatura) pgx.CopyFromSource {
	return &faturaSource{faturas: faturas}
}

// Utilizado pelo CopyFrom para gerenciar os dados
func (fs *faturaSource) Next() bool {
	fs.pos++
	return fs.pos <= len(fs.faturas)
}

// Utilizado pelo CopyFrom para gerenciar os dados
func (fs *faturaSource) Values() ([]any, error) {
	fatura := fs.faturas[fs.pos-1]

	return []any{
		fatura.Emitente, fatura.Documento, fatura.Contrato, fatura.Categoria, fatura.QtdNota, fatura.Fatura, fatura.Valor, fatura.DataCompra, fatura.DataPagamento,
	}, nil
}

func (fs *faturaSource) Err() error {
	return nil
}
