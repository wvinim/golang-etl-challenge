#!/bin/bash
REPS=10
STR="emitente	documento	contrato	categoria	qtdNota	fatura	valor	data_compra	data_pagamento
BlaBlabla Soluções Digitais	401678cc1d9a716baebbac87452f62686572f729677c24921b3f08925d468e4c	17/2024	Prestação de serviços	4	34239	20.202.953,58	26/11/2024	30/01/2025
BlaBlabla Consultoria		26/2024	Prestação de serviços	2	14471	5.446.671,88	26-01-2025	25-03-2025
BlaBlabla Tech Group	8e0d270e4169d611a6d3fa6f3a8f266c283b86f16385ece316ebbee0aeba2446	21/2024	Prestação de serviços	4	66649	11.406.536,52	05/05/2025	19/05/2025
BlaBlabla Estratégias		19/2024	Locações	4	67792	1.463.301,36	01-12-2024	18-02-2025
BlaBlabla Services		08/2025	Locações	1	95229	29.112.385,18	22-03-2025	07-04-2025"
FILE="base_ficticia_$REPS.txt"
> "$FILE"
for ((i=1;i<=REPS;i++)); do
    echo "$STR" >> "$FILE"
done
echo "Arquivo '$FILE' criado com sucesso!"