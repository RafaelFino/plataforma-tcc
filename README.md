# plataforma-tcc
Guia para o desenvolvimento do projeto de conclusão de curso

## Objetivos
Colocar em prática todos os conceitos aprendidos durante todo o curso preparatório da Plataforma Impact, criando um projeto rico, seguindo uma especificação aderente ao mundo real de uma empresa para ser o primeiro grande item de portfolio no curriculo dos estudantes

## Instruções
- Liberdade de escolha das tecnologias, desde que os requisitos sejam planamente atendidos
- O estudante deve criar um repositório no github.com com o nome **plataforma-tcc**

## Requisitos de negócio
Criar um sistema composto por serviços independentes:
- Cadastro de clientes
- Cadastro de produtos
- Serviço de cotações de moedas estrangeiras
- Controle de vendas

### Cadastro de clientes
Esse serviço deve ser capaz de:
- Cadastrar um cliente
- Alterar um cliente previamente cadastrado
- Inativar um cliente (delete lógico)
- Listar todos os clientes

#### Modelo de cliente
| Nome | Tipo | Descrição | 
|:-|:-|:-| 
| __id__ | STRING | Identificador único do cliente |
| __name__ | STRING | Nome do cliente | 
| __surname__ | STRING | Sobrenome do cliente |
| __created_at__ | TIMESTAMP | Data de criação do registro | 
| __updated_at__ | TIMESTAMP | Data de criação do registro | 


### Cadastro de produtos
Esse serviço deve ser capaz de:
- Cadastrar um produto
- Alterar um produto previamente cadastrado
    - Alterar suas caracteristicas
    - Alterar suas quantidades em estoque
    - Alterar seu preço em BRL
- Inativar um produto
- Listar todos os produtos, com o preço em todas as moedas disponíveis no serviço de cotações
- Listar os produtos, seguindo filtros de:
    - Preço por moeda (maior ou igual, menor ou igual)
    - Nome ou descrição
    - Quantidade em estoque

#### Modelo de produto
| Nome | Tipo | Descrição | 
|:-|:-|:-| 
| __id__ | STRING | Identificador único do produto |
| __name__ | STRING | Nome do produto | 
| __dest__ | STRING | Descrição do produto | 
| __quantity__ | INTEGER | Quantidade em estoque desse produto |
| __created_at__ | TIMESTAMP | Data de criação do registro | 
| __updated_at__ | TIMESTAMP | Data de criação do registro | 

### Serviço de cotações de moedas estrangeiras
Esse serviço deve ser capaz de fornecer a cotação das seguintes moedas:
- __BRL__ vs __EUR__
- __BRL__ vs __USD__
- __BRL__ vs __GPB__
- __BRL__ vs __CNY__

Os códigos de moeda, seguem a [ISO-4717](https://pt.iban.com/currency-codes)

Esse serviço deve ser capaz de informar a cotação de uma moeda em específico por meio de rota HTTP e também informar todas elas caso nenhuma moeda seja especificada

Esse serviço pode consultar uma API externa para conseguir os dados [API de cotações](https://economia.awesomeapi.com.br/all), porém essa API possui limitação de requisições, portanto é necessário um mecanismo de cache para que essa consulta aconteça no máximo uma vez ao dia e a data de atualização do preço deve ser informada sempre que os valores forem consultados

#### Modelo de cotações
| Nome | Tipo | Descrição | 
|:-|:-|:-| 
| __code__ | STRING | Código da moeda | 
| __value__ | FLOAT | Descrição do produto | 
| __created_at__ | TIMESTAMP | Data de criação do registro | 

### Controle de vendas
Esse serviço deve ser capaz de se comunicar com os serviços de:
- Clientes
- Produtos
- Cotações

Funcionalidades:
- Criar uma venda para um cliente, puxando os dados do cliente via API
- Incluir produtos nessa venda, verificando se existe essa quantidade disponível no serviço de produtos, via API
- Alterar produtos nessa venda (quantidade), verificando a disponibilidade, como na inclusão, via API
- Efetivar uma venda previamente criada e efetivando o uso do produto, ou seja, reduzindo seu estoque disponível no serviço de produtos, via API
- No momento da efetivação da venda, deve ser informado o preço total em todas as moedas disponíveis no serviço de cotações
- Consultar todas as vendas de um determinado produto
- Consultar o total de vendas, por estado, de um determinado produto
- Consultar todas as vendas, por estado
- Cancelar uma determinada venda, previamente criada

#### Modelo de venda
| Nome | Tipo | Descrição |
|:-|:-|:-| 
| __id__ | STRING | Identificador único de cada venda |
| __client_id__ | STRING | Identificador do cliente da venda |
| __array<item_venda>__ | ARRAY | Itens contidos na venda |
| __status__ | INTEGER | Identificador do estado da venda |
| __created_at__ | TIMESTAMP | Data de criação do registro | 
| __updated_at__ | TIMESTAMP | Data de criação do registro | 

#### Modelo de item de venda
| Nome | Tipo | Descrição |
|:-|:-|:-| 
| __sell_id__ | STRING | Identificador único de cada venda, referência a venda a qual esse item pertence |
| __product_id__ | STRING | Identificador do produto contido na venda | 
| __quantity__ | INTEGER | Quantidade do produto nessa venda |
| __created_at__ | TIMESTAMP | Data de criação do registro | 
| __updated_at__ | TIMESTAMP | Data de criação do registro |

#### Estados possíveis de uma venda (INTEGER)
| Valor | Tipo | Descrição |
|:- |:- | :- | 
| __0__| __STARTED__ | Venda criada |
| __1__| __PROGRESS__ | Venda em progresso (quando o primeiro item é adicionado) |
| __2__| __DONE__ | Venda finalizada, o estoque deve ser validado e se possível alterado |
| __3__ | __CANCELED__ | Venda cancelada |


## Requisitos técnicos
- Todas as APIs devem ser HTTP REST
- Todos os serviços devem possuir logs, conforme padrão descrito nesse documento:
    - Logs para todos os métodos
    - Logs para todas as transações que alterem algum dado ou estado
    - Logs para todas as inicializações ou paradas de serviços
    - Logs de tempo para cada requisição em milisegundos (tempo gasto entre a recepção da requisição e a resposta)
- Todas as APIs devem possuir testes unitários, com liberdade de escolha da ferramenta (postman, curl via script e etc)
- Todos os métodos devem ser testáveis e devem possuir testes
- Todos os serviços devem estar dentro de um único arquivo docker-compose, capaz de subir toda a infra necessária para o funcionamento de toda a plataforma, com cada serviço rodando necessariamente em um container independente
- O banco de dados deve ser específico para cada serviço: cada serviço tem o seu repositório, podendo assim ser inclusive, se o aluno optar, em diferentes tecnologias. Não será permitido que um serviço tenha acesso ao banco de dados do outro serviço
- Cada serviço pode ser escrito em diferentes tecnologias, desde que respeitem o padrão de log e todos os requisitos já listados
- A solução deve seguir um modelo de arquitura de solução com segregação de responsabilidades:
    - Camadas de domínios/entidades
    - Camadas de serviços/controler
    - Camadas de repositórios/storage
    - Camadas de APIs/Handlers
- Todas as respostas dos serviços devem ser no formato JSON e ter obrigatoriamente os seguintes campos:

| Campo | Tipo | Descrição |
| :- | :- | :- |
| __message__ | STRING | Mensagem amigável de retorno |
| __timestamp__ | TIMESTAMP | Momento em que o servidor retornou | 
| __elapsed__ | INTEGER | Tempo, em milisegundos, que o servidor demorou para processar a requisição | 
| __error__ | STRING | Campo para descrição do erro, só deve ser retornado em caso de erro no servidor | 

### Modelo de log
Cada campo deve ser separado por um TAB (__\t__) e finalizado com um LINEFEED (__\n__) (ver descrição de [ASCII](https://pt.wikipedia.org/wiki/ASCII))

| Campo | Descrição | 
| :- | :- | 
| __timestamp__ | Momento em que a entrada foi gerada, deve seguir a [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) | 
| __level__ | Nível de criticidade da entrada | 
| __message__ | Mensagem de log | 

### Nível de log
| Nível | Descrição | 
| :- | :- | 
| __DEBUG__ | Nível de criticidade de depuração, não deve estar ativo em ambiente produtivo (no momento da avaliação) | 
| __INFO__ | Nível de criticidade de informação, deve seguir a especificação de logs | 
| __WARNING__ | Nível de criticidade da alerta, algo aconteceu fora do esperado, mas a rotina deve continuar | 
| __ERROR__ | Nível de erro, a aplicação não é capaz de processar o pedido e deve sinalizar que um erro aconteceu | 
| __CRITICAL__ | Nível de erro crítico, a aplicação não é capaz de continuar e será encerrada | 

```
2024-03-18T23:10:57Z    INFO    Service starting...
2024-03-18T23:10:57Z    INFO    Connecting on database
2024-03-18T23:10:57Z    DEBUG    Database running under localhost:5432
2024-03-18T23:10:58Z    INFO    Listening port 8000
```