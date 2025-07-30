package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	websocketHandler = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ServerWs(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketHandler.Upgrade(w, r, nil)
	if err != nil {
		log.Printf(`unable to upgrade %v`, err)
		return
	}
	log.Printf(`upgraded %v`, conn)

}