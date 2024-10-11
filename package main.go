package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Permitir conexões de qualquer origem, por questões de segurança, ajuste conforme necessário
    },
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    // Faz o upgrade da conexão HTTP para uma conexão WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    log.Println("Cliente conectado")

    for {
        // Lê mensagem do cliente
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        // Log da mensagem recebida
        log.Printf("Mensagem recebida: %s", p)

        // Envia a mensagem de volta ao cliente (echo)
        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }
    }
}

func setupRoutes() {
    http.HandleFunc("/ws", wsEndpoint)
}

func main() {
    log.Println("Iniciando servidor WebSocket na porta 8080")
    setupRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
