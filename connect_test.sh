#!/bin/bash

# URL do token
TOKEN_URL="http://172.17.0.1:8081/get-token"
# URL base do WebSocket
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

# Nomes dos canais
CANAIS=("geral" "canal1")
# Número de mensagens a serem enviadas por canal
NUM_MENSAGENS=5

# Enviar mensagens "oi" em loop para cada canal
for CANAL in "${CANAIS[@]}"; do
    for ((i=1; i<=NUM_MENSAGENS; i++)); do
        {
            # Usar websocat para se conectar e enviar mensagens com o canal dentro da mensagem
            FULL_WS_URL="$WS_URL?token=$TOKEN"
            echo "Conexão para o WebSocket estabelecida para o canal $CANAL"
            MENSAGEM="{\"canal\": \"$CANAL\", \"mensagem\": \"oi $i\"}"
            echo "Enviando mensagem: $CANAL $MENSAGEM"
            echo "$MENSAGEM" | websocat "$FULL_WS_URL" -E  # -E força websocat a fechar a conexão
            sleep 2  # Intervalo de 2 segundos entre as mensagens
        } &
    done
done

# Espera indefinidamente até que o sinal SIGINT seja recebido
wait
