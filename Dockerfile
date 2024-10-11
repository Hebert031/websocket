# Use uma imagem base do Go
FROM golang:1.20-alpine

# Defina o diretório de trabalho no contêiner
WORKDIR /app

# Copie os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixe as dependências
RUN go mod tidy

# Copie o restante dos arquivos do projeto
COPY . .

# Compile o aplicativo
RUN go build -o websocket-server websocket.go

# Defina o comando para executar o aplicativo
CMD ["./websocket-server"]

# Exponha a porta que o aplicativo vai usar
EXPOSE 8080
