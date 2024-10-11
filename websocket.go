package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var jwtSecret = []byte("CampTecnologia@2024!*1987!1988//")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Estrutura para armazenar os clientes e seus canais
var channels = make(map[string]map[*websocket.Conn]bool)

func createJWT() (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"user":       "user1",
		"exp":        time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de assinatura inválido")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func broadcastMessage(channel, message string) {
	// Enviar a mensagem para todos os clientes do canal
	for client := range channels[channel] {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			client.Close()
			delete(channels[channel], client)
		}
	}
}

func broadcastClients(channel string) {
	clientCount := len(channels[channel])
	clientList := map[string]int{"clientesConectados": clientCount}
	jsonData, err := json.Marshal(clientList)
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}
	broadcastMessage(channel, string(jsonData))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	channel := r.URL.Query().Get("channel") // Nome do canal vindo da URL

	if tokenString == "" {
		http.Error(w, "Autorização requerida", http.StatusUnauthorized)
		return
	}

	token, err := validateToken(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao atualizar para WebSocket:", err)
		return
	}
	defer ws.Close()

	// Verificar se o canal já existe, se não, criar um novo
	if channels[channel] == nil {
		channels[channel] = make(map[*websocket.Conn]bool)
	}

	// Adicionar o cliente ao canal
	channels[channel][ws] = true
	broadcastClients(channel)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler a mensagem:", err)
			delete(channels[channel], ws)
			broadcastClients(channel)
			break
		}

		if len(msg) > 1024 {
			fmt.Println("Mensagem muito longa, descartando.")
			continue
		}

		fmt.Printf("Mensagem recebida no canal %s: %s\n", channel, msg)
		broadcastMessage(channel, string(msg))
	}
}

func main() {
	http.HandleFunc("/get-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := createJWT()
		if err != nil {
			http.Error(w, "Erro ao criar token", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, token)
	})

	http.HandleFunc("/ws", handleConnections)

	corsHandler := corsMiddleware(http.DefaultServeMux)

	fmt.Println("WebSocket iniciado na porta :8080")
	err := http.ListenAndServe(":8080", corsHandler)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
