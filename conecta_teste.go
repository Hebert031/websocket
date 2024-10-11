package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Função para buscar o token
func getToken(tokenURL string) (string, error) {
	resp, err := http.Get(tokenURL)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar token: %v", err)
	}
	defer resp.Body.Close()

	var token string
	_, err = fmt.Fscan(resp.Body, &token)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o token: %v", err)
	}
	return token, nil
}

// Função para enviar mensagem para o WebSocket
func sendMessage(wsURL, token, canal, mensagem string) error {
	// Conectar ao WebSocket
	fullURL := fmt.Sprintf("%s?token=%s", wsURL, token)
	conn, _, err := websocket.DefaultDialer.Dial(fullURL, nil)
	if err != nil {
		return fmt.Errorf("erro ao conectar no WebSocket: %v", err)
	}
	defer conn.Close()

	// Enviar a mensagem
	msg := fmt.Sprintf(`{"canal": "%s", "mensagem": "%s"}`, canal, mensagem)
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("erro ao enviar mensagem: %v", err)
	}

	fmt.Printf("Mensagem enviada para o canal %s: %s\n", canal, mensagem)
	return nil
}

func main() {
	// URLs do token e do WebSocket
	tokenURL := "http://172.17.0.1:8081/get-token"
	wsURL := "ws://172.17.0.1:8081/ws"

	// Buscar o token
	token, err := getToken(tokenURL)
	if err != nil {
		log.Fatalf("Erro ao obter o token: %v\n", err)
	}
	fmt.Printf("Token obtido: %s\n", token)

	// Definir canais e número de mensagens
	canais := []string{"geral", "canal1"}
	numMensagens := 5

	// Enviar mensagens para os canais
	for _, canal := range canais {
		for i := 1; i <= numMensagens; i++ {
			mensagem := fmt.Sprintf("oi %d", i)
			err := sendMessage(wsURL, token, canal, mensagem)
			if err != nil {
				log.Printf("Erro ao enviar mensagem para o canal %s: %v\n", canal, err)
			}
			time.Sleep(2 * time.Second) // Esperar 2 segundos entre as mensagens
		}
	}

	// Finalizar o programa
	fmt.Println("Mensagens enviadas com sucesso.")
}
