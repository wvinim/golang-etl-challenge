package normalizer

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// Tabela de busca para normalização de valores float
var tableStringToFloat [256]byte

// Tabela de busca para normalização de caracteres
var tableCharacters [256]byte

// Preenche as tabelas de normalização quando o programa inicia
func init() {
	// Inicializa a tableStringToFloat com 0 para chars ignorados
	for i := range tableStringToFloat {
		tableStringToFloat[i] = 0
	}

	// Mapeia os dígitos
	for i := '0'; i <= '9'; i++ {
		tableStringToFloat[i] = byte(i)
	}

	// Mapeia a vírgula para um ponto
	tableStringToFloat[','] = '.'

	// ----------------------------------

	// Inicializa a tableStringToFloat com 0 para chars ignorados
	for i := range tableCharacters {
		tableCharacters[i] = 0
	}

	// Mapeia as letras minúsculas ASCII para maiúsculas
	for c := 'a'; c <= 'z'; c++ {
		tableCharacters[c] = byte(c - 32)
	}

	// Mapeia as letras maiúsculas
	for c := 'A'; c <= 'Z'; c++ {
		tableCharacters[c] = byte(c)
	}

	// Mapeia os números
	for c := '0'; c <= '9'; c++ {
		tableCharacters[c] = byte(c)
	}

	// Mapeia caracteres específicos
	tableCharacters['-'] = '-'
	tableCharacters['.'] = '.'
	tableCharacters[' '] = ' '

	// Mapeia os acentos mais comuns PT-BR
	for _, c := range "áàãâäÁÀÃÂÄ" {
		tableCharacters[c] = 'A'
	}
	for _, c := range "éèêëÉÈÊË" {
		tableCharacters[c] = 'E'
	}
	for _, c := range "íìîïÍÌÎÏ" {
		tableCharacters[c] = 'I'
	}
	for _, c := range "óòõôöÓÒÕÔÖ" {
		tableCharacters[c] = 'O'
	}
	for _, c := range "úùûüÚÙÛÜ" {
		tableCharacters[c] = 'U'
	}
	tableCharacters['ç'], tableCharacters['Ç'] = 'C', 'C'
}

// Limpa a string de um número usando a tabela de busca
func NormalizeStringToFloat(s string) float64 {
	// O strings.Builder é usado para construir a string final, minimizando alocações de memória.
	var builder strings.Builder
	// A função Grow é usada para pré-alocar a capacidade necessária, otimizando a performance.
	builder.Grow(len(s))

	// Percorre os caracteres da string e busca eles na tabela fixa para substituição
	for i := 0; i < len(s); i++ {
		b := tableStringToFloat[s[i]]
		if b != 0 {
			builder.WriteByte(b)
		}
	}

	// Realiza o parse para o tipo float64 e retorna o valor
	valor, _ := strconv.ParseFloat(builder.String(), 64)
	return valor
}

// Normaliza a string utilizando uma estratégia híbrida de decodificação e substituição
func NormalizeCharactersHybrid(s string) *string {
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
		if out := tableCharacters[charIndex]; out != 0 {
			b.WriteByte(out)
		}
	}
	return StringOrNil(b.String())
}

// Otimiza a conversão de datas de DD/MM/AAAA ou DD-MM-AAAA
// para YYYY-MM-DD
func NormalizeDateString(s string) *string {
	if len(s) != 10 {
		return nil // Retorna nil para formatos inválidos
	}

	// Extrai os componentes da data com base nas posições fixas
	dia := s[0:2]
	mes := s[3:5]
	ano := s[6:10]

	// Constrói a string no novo formato YYYY-MM-DD
	// Usando um array de bytes para evitar múltiplas alocações
	var result [10]byte
	copy(result[0:4], ano)
	result[4] = '-'
	copy(result[5:7], mes)
	result[7] = '-'
	copy(result[8:10], dia)

	return StringOrNil(string(result[:]))
}

// Retorna nil caso a string seja vazia
func StringOrNil(s string) *string {
	if s == "" || s == "-" {
		return nil
	}
	return &s
}
