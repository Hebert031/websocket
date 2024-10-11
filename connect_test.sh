#!/bin/bash

# URL do token
TOKEN_URL="http://172.17.0.1:8081/get-token"
# URL do WebSocket
WS_URL="ws://172.17.0.1:8081/ws"

# Função para encerrar conexões
cleanup() {
    echo "Encerrando conexões..."
    pkill -f "websocat $WS_URL"
    exit 0
}

# Captura o sinal SIGINT (Ctrl + C)
trap cleanup SIGINT

# Buscar o token
TOKEN=$(curl -s $TOKEN_URL)

if [ -z "$TOKEN" ]; then
    echo "Erro ao obter o token."
    exit 1
fi

echo "Token obtido: $TOKEN"

# Número de conexões
NUM_CONNECTIONS=5

# Enviar mensagens "oi" em loop
for ((i=1; i<=NUM_CONNECTIONS; i++)); do
    {
        # Usar websocat para se conectar e enviar mensagens
        echo "Conexão $i estabelecida"
        while true; do
            echo "oi" | websocat "$WS_URL?token=$TOKEN"
            sleep 2  # Intervalo de 2 segundos entre as mensagens
        done &
    }
done

# Espera indefinidamente até que o sinal SIGINT seja recebido
wait
