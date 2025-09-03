## Lista de apoio aos estudos, comandos úteis e pesquisas

### Neste arquivo eu vou adicionar algumas descobertas em geral que me ajudaram durante os estudos

### Para iniciar um projeto em go

    go mod init golang-etl-challenge

## Para adicionar pacote externo (ex golang.org/x/text)

    go get golang.org/x/text

## Para analisar a performance é preciso:

    - Criar testes para realizar o benchmark
    - Criar um output dos profiles da CPU e memória
    - Analisar com a ferramente pprof
        - Para gerar svg e usar o comando web é preciso ter o pacote graphviz (brew install graphviz)

## Para criar profile de cpu

    go test -bench=. -cpuprofile=cpu.out

## Para criar profile de memória

    go test -bench=. -memprofile=mem.out

## Para analisar o profile da cpu em modo web

    go tool pprof -http=:8080 cpu.out

##  Para analisar o profile da cpu em modo web com foco em uma função específica

    go tool pprof -http=:8080 -focus=normalizeFast cpu.out