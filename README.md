# 🚀 Desafio ETL em GoLang

![Badge de Status](https://img.shields.io/badge/Status-Conclu%C3%ADdo-green)
![Badge de Licença](https://img.shields.io/github/license/wvinim/golang-etl-challenge)

## 📝 Descrição do Projeto

Solução para o desafio de criar um serviço para o carregamento de dados em lote.

## 🌟 Funcionalidades Principais

* 🎯 Funcionalidade 1: Recebimento do arquivo bruto através de um serviço HTTP.
* 📦 Funcionalidade 2: Leitura, processamento dos dados e persistência no banco com paralelismo.
* ⚡ Funcionalidade 3: Foco em performance extrema e estabilidade para arquivos com grandes volumes.

## 💻 Tecnologias Utilizadas

* **Backend:** `GoLang` (imagem 1.25-alpine)
* **Banco de Dados:** `PostgreSQL` (imagem 15-alpine)
* **Outros:** `Docker` (versão 28.3.0)

## ⚙️ Antes de rodar

Crie um arquivo `.env` na raiz do projeto, use como exemplo o `.env.example`.

Caso decida rodar sem o docker ou tiver algum problema com as portas, veja a legenda das variáveis:
```bash
POSTGRES_HOST = Endereço de rede do seu servidor PostgreSQL
POSTGRES_USER = Usuário do seu servidor PostgreSQL
POSTGRES_PASSWORD = Senha de acesso do seu usuário ao seu servidor PostgreSQL
POSTGRES_DB = Nome da database dentro do seu servidor PostgreSQL
POSTGRES_PORT = Porta de rede utilizada pelo seu servidor PostgreSQL
CHUNK_SIZE = Número de registros a ser carregado por lote
ENVIRONMENT = Ambiente em que este aplicatibo está rodando
SERVER_PORT = Porta de rede em que o servidor HTTP estará disponível
```

## ⚙️ Como Rodar em ambiente local

Com o GoLang (1.25) instalado em sua máquina, siga os passos abaixo:

1.  Clone este repositório:
    ```bash
    git clone [https://github.com/wvinim/golang-etl-challenge.git](https://github.com/wvinim/golang-etl-challenge.git)
    ```
2.  Entre na pasta do projeto:
    ```bash
    cd golang-etl-challenge
    ```
3.  Configure o seu arquivo `.env` com as credenciais do seu PostgreSQL

4.  Instale as dependências:
    ```bash
    go mod init golang-etl-challenge
    ```
5.  Com o comando abaixo, seu serviço estará pronto para receber requisições:
    ```bash
    go run main.go
    ```

## ⚙️ Como Rodar em Docker

Com o docker instalado em sua máquina, siga os passos abaixo:

1.  Clone este repositório:
    ```bash
    git clone [https://github.com/wvinim/golang-etl-challenge.git](https://github.com/wvinim/golang-etl-challenge.git)
    ```
2.  Entre na pasta do projeto:
    ```bash
    cd golang-etl-challenge
    ```
3.  Verifique o seu arquivo `.env`

4.  Com o comando abaixo, seu serviço estará pronto para receber requisições:
    ```bash
    docker-compose up -d
    ```

## ⚙️ Como Remover este aplicativo do seu docker

```bash
docker-compose down --rmi all
```

## 🗺️ Como Usar

Seja via local ou docker, você precisará de cliente http para enviar seu arquivo, segue o cURL de exemplo:

```bash
curl --request POST \
  --url http://localhost:8088/upload \
  --header 'Content-Type: multipart/form-data' \
  --form file=@/caminho/para/seu/arquivo/base_ficticia_dados_prova.txt
```

Para gerar um arquivo fictício, use o script `generate.sh`, altere a variável REPS para a quantidade de registros desejada:
```bash
chmod +x generate.sh

./generate.sh
```

## 📝 Documentação da tabela `faturas`

Segue o script e considerações:

```bash
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
```

    - O próprio script cria a tabela durante a inicialização (se não existir), não utilizei gORM devido a simplicidade do projeto.
    - Não utilizei valores fixos para o tamanho dos campos (ex varchar(50)), por dois motivos:
        - A decisão do tamanho apenas pela amostragem da base fictícia pode trazer problemas futuros em um caso real
        - A diferença de performance é mínima, mas é cabível caso exista um contrato definido para este arquivo

Exemplos de  índices:

```bash
CREATE INDEX IF NOT EXISTS idx_faturas_emitente ON faturas (emitente);
CREATE INDEX IF NOT EXISTS idx_faturas_data_compra ON faturas (data_compra);
CREATE INDEX IF NOT EXISTS idx_faturas_emitente_data_compra ON faturas (emitente, data_compra);
```

    - O uso de índices é crucial para otimizar as consultas em tabelas realmente grandes
        - Porém, comprometem um pouco a performance do insert
        
## 💻 Diário de bordo

Criei estes arquivos de apoio para desenvolver a estratégia e organizar minha aprendizagem em Go:

* **[TODO.md](TODO.md):** Lista de passos que trilhei para a conclusão do desafio
* **[LEARNING.md](LEARNING.md):** Lista de apoio aos estudos, comandos úteis e pesquisas
* **[CHALLENGE.md](CHALLENGE.md):** Lista de requisitos e critérios de avaliação
* **Comentários no código:** Deixei o código extremamente comentado para justificar as decisões

## 🌟 Considerações finais

- Já fiz algo parecido em Python então conhecia um pouco o "caminho das pedras"
- Foi meu primeiro desafio relevante em Go
- Senti o mesmo entusiasmo de 2001, quando comecei em C, com o plus da senioridade e mais recursos
- É extremamente gratificante desenvolver esse tipo de solução em Go

## 🤝 Como Contribuir

Contribuições são sempre bem-vindas! Se você quiser ajudar, por favor, siga os passos abaixo:

1.  Crie um fork do projeto.
2.  Crie uma nova branch com a sua feature: `git checkout -b minha-nova-feature`.
3.  Faça o commit das suas alterações: `git commit -m 'feat: minha nova feature'`.
4.  Envie para a branch: `git push origin minha-nova-feature`.
5.  Abra uma Pull Request.

## 📄 Licença

Este projeto está sob a licença [MIT](LICENSE) - veja o arquivo `LICENSE` para mais detalhes.

## 👤 Autor

* **Wellington Moraes** - [GitHub](https://github.com/wvinim) | [LinkedIn](https://www.linkedin.com/in/wellington-vinicius-moraes-726b4b58/)

---

> ✨ Se você gostou do projeto, por favor, dê uma estrela! Obrigado! ✨