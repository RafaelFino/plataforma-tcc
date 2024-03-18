# plataforma-tcc
Guia para o desenvolvimento do projeto de conclusão de curso

## Objetivos
Colocar em prática todos os conceitos aprendidos durante todo o curso preparatório da Plataforma Impact, criando um projeto rico, seguindo uma especificação aderente ao mundo real de uma empresa para ser o primeiro grande item de portfolio no curriculo dos estudantes

## Instruções
- Liberdade de escolha das tecnologias, desde que os requisitos sejam planamente atendidos
- O estudante deve criar um repositório no github.com com o nome **plataforma-tcc**

## Requisitos
Criar um sistema composto por serviços independentes:
- Cadastro de clientes
- Cadastro de produtos
- Serviço de cotações de moedas estrangeiras
- Controle de vendas

### Requisitos técnicos
- Todas as APIs devem ser HTTP REST
- Todos os serviços devem possuir logs, conforme padrão descrito nesse documento:
    - Logs para todos os métodos
    - Logs para todas as transações que alterem algum dado ou estado
    - Logs para todas as inicializações ou paradas de serviços
    - Logs de tempo para cada requisição em milisegundos (tempo gasto entre a recepção da requisição e a resposta)
- Todas as APIs devem possuir testes unitários, com liberdade de escolha da ferramenta (postman, curl via script e etc)
- Todos os serviços devem estar dentro de um único arquivo docker-compose, capaz de subir toda a infra necessária para o funcionamento de toda a plataforma, com cada serviço rodando necessariamente em um container independente
- O banco de dados pode ser específico para cada serviço: cada serviço tem o seu repositório, podendo assim ser inclusive, se o aluno optar, em diferentes tecnologias
- Cada serviço pode ser escrito em diferentes tecnologias, desde que respeitem o padrão de log e todos os requisitos já listados

#### Modelo de log
| Campo | Descrição | 
| :- | :- | 
| __timestamp__ | Momento em que a entrada foi gerada | 
| __level__ | Nível de criticidade da entrada | 
| __message__ | Mensagem de log | 

#### Nível de log
| Nível | Descrição | 
| :- | :- | 
| __DEBUG__ | Nível de criticidade de depuração, não deve estar ativo em ambiente produtivo (no momento da avaliação) | 
| __INFO__ | Nível de criticidade de informação, deve seguir a especificação de logs | 
| __WARNING__ | Nível de criticidade da alerta, algo aconteceu fora do esperado, mas a rotina deve continuar | 
| __ERROR__ | Nível de erro, a aplicação não é capaz de processar o pedido e deve sinalizar que um erro aconteceu | 
| __CRITICAL__ | Nível de erro crítico, a aplicação não é capaz de continuar e será encerrada | 

### Cadastro de clientes
Esse serviço deve ser capaz de:
- Cadastrar um cliente
- Alterar um cliente previamente cadastrado
- Inativar um cliente (delete lógico)
- Listar todos os clientes

#### Modelo de cliente
| Nome | Tipo | Descrição | 
|:-|:-:|:-| 
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
|:-|:-:|:-| 
| __id__ | STRING | Identificador único do produto |
| __name__ | STRING | Nome do produto | 
| __dest__ | STRING | Descrição do produto | 
| __quantity__ | float | Quantidade em estoque desse produto |
| __created_at__ | TIMESTAMP | Data de criação do registro | 
| __updated_at__ | TIMESTAMP | Data de criação do registro | 

### Serviço de cotações de moedas estrangeiras
Esse serviço deve ser capaz de fornecer a cotação das seguintes moedas:
- BRL vs EUR
- BRL vs USD
- BRL vs GPB
- BRL vs CNY

Os códigos de moeda, seguem a [ISO-4717](https://pt.iban.com/currency-codes)

Esse serviço deve ser capaz de informar a cotação de uma moeda em específico por meio de rota HTTP e também informar todas elas caso nenhuma moeda seja especificada

Esse serviço pode consultar uma API externa para conseguir os dados [API de cotações](https://economia.awesomeapi.com.br/all), porém essa API possui limitação de requisições, portanto é necessário um mecanismo de cache para que essa consulta aconteça no máximo uma vez ao dia.

#### Modelo de cotações
| Nome | Tipo | Descrição | 
|:-|:-:|:-| 
| __code__ | STRING | Código da moeda | 
| __value__ | float | Descrição do produto | 
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

## Arquitetura