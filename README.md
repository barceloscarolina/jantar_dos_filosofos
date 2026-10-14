Jantar dos Filósofos

Projeto desenvolvido em Go para a disciplina de Fundamentos de Processamento Paralelo e Distribuído (PUC-RS), comparando quatro versões do problema dos filósofos jantando.

V0: execução sequencial.

V1: execução concorrente sujeita a deadlock.

V2: prevenção de deadlock com sincronização excessiva.

V3: prevenção de deadlock por ordenação dos recursos.

Execução:

Execute cada versão na raiz do projeto:

go run V0.go
go run V1.go
go run V2.go
go run V3.go

Cada filósofo executa 20 ciclos, com tempos simulados de 2 segundos para pensar e 3 segundos para comer.

As versões que concluem normalmente são executadas cinco vezes para coleta e comparação das métricas de desempenho.

Linguagem: Go 1.27.1.
