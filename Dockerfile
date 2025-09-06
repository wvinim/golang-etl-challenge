# Estágio de compilação
# Use uma imagem com o ambiente de Go para compilar.
FROM golang:1.25-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum para o cache de dependências
COPY go.mod ./
COPY go.sum ./

# Baixa as dependências
RUN go mod download

# Copia todo o código-fonte
COPY . .

# Compila a aplicação e cria um binário estático.
# A flag -o define o nome do arquivo de saída.
# A flag -ldflags "-s -w" reduz o tamanho do binário removendo tabelas de símbolos e informações de depuração.
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Estágio final (imagem de execução)
# Use uma imagem extremamente leve (alpine) para a execução.
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /root/

# Copia o binário compilado do estágio anterior
COPY --from=builder /app/main .

# Comando para rodar a aplicação
CMD ["./main"]