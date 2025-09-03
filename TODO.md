## Lista de passos que considero essenciais para a conclusão do desafio

- [x] Como medir a performance em golang? 
    - Aprendi a gerar os benchmarks e analisar os principais indicadores de performance (detalhes sobre os comandos em LEARNING.md)
    - Deixei o exemplo prático utilizado na pasta steps/1-performance (O exemplo serviu para o próximo passo)
- [x] Qual a estratégia mais rápida para normalização das strings?
    - De/para com um range fixado de caracteres conhecidos (pt-BR)
    - Decodificação híbrida para diferenciar caracteres ASCII dos demais.
    - O exemplo prático do melhor resultado que obtive (até agora rss) está em steps/2-normalize
        - 9120418 operações em 2.563s, 131.3 ns/op
- [ ] Como ler grandes arquivos em buffer?
    - Faz diferença transformar esse arquivo em binário puro antes?
- [ ] Como enviar inserts em batch para o banco?