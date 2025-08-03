# WebSocket com JWT

Este projeto implementa um servidor WebSocket usando a linguagem de programação Go. O servidor oferece autenticação baseada em JWT (JSON Web Token) para garantir que apenas usuários autorizados possam se conectar ao WebSocket.

## Funcionalidades

- **Criação de Token JWT**: Endpoint para gerar um token JWT.
- **Validação de Token JWT**: Verifica o token JWT enviado na URL para autenticação.
- **WebSocket**: Conexão WebSocket que ecoa as mensagens enviadas pelo cliente.
- **CORS**: Suporte a CORS para permitir conexões de origens diferentes.

## Dependências

- [Golang](https://golang.org/dl/) (você deve ter o Go instalado)
- [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket)
- [Golang JWT](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)

## Configuração

1. **Clone o Repositório**

   ```bash
   git clone <URL_DO_REPOSITORIO>
   cd <NOME_DO_REPOSITORIO>
Instale as Dependências

bash
Copiar código
go mod tidy
Execute o Servidor

bash
Copiar código
go run main.go
O servidor irá iniciar na porta 8080.

Endpoints
GET /get-token

Gera um token JWT. Utilize este endpoint para obter um token necessário para se conectar ao WebSocket.

bash
Copiar código
curl http://localhost:8080/get-token
GET /ws?token=<TOKEN>

Conecta ao WebSocket protegido. Substitua <TOKEN> pelo token JWT obtido do endpoint /get-token.

Exemplo de conexão com websocat:

bash
Copiar código
websocat "ws://localhost:8080/ws?token=<TOKEN>"
Middleware CORS
O servidor inclui um middleware para lidar com CORS, permitindo requisições de qualquer origem e suportando os métodos GET, POST e OPTIONS.

Estrutura do Código
main.go: Arquivo principal que contém a lógica do servidor, criação e validação do token JWT, e configuração do WebSocket.

createJWT(): Gera um novo token JWT.
validateToken(tokenString string): Valida um token JWT.
corsMiddleware(next http.Handler): Middleware para adicionar cabeçalhos CORS.
handleConnections(w http.ResponseWriter, r *http.Request): Lida com as conexões WebSocket.
main(): Configura os endpoints e inicia o servidor HTTP.

Segurança
A chave secreta para assinatura do JWT está definida como uma variável global jwtSecret no código. Para produção, considere usar uma chave mais segura e armazenar a chave secreta em um ambiente seguro.


Contato
Para mais informações, entre em contato.
