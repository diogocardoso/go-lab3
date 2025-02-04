# GO Expert - lab 3
Projeto do Laboratório "Concorrência com Golang - Leilão" do treinamento GoExpert(FullCycle).

## O desafio

Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, [o projeto clonado](https://github.com/devfullcycle/labs-auction-goexpert) necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.


## Como rodar o projet
``` shell
## put the docker-compose containers up
docker-compose up
ou
docker-compose up --build
```

## Como parar os containers
``` shell
docker-compose down

```

## Requisitos: implementação
- Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- Um teste para validar se o fechamento está acontecendo de forma automatizada;

## Requisitos: entrega
- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.
