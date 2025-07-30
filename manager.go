package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketHandler = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients map[*Client]bool
	mu sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[*Client]bool),
		mu: sync.Mutex{},
	}
}

func (m *Manager) ServerWs(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketHandler.Upgrade(w, r, nil)
	if err != nil {
		log.Printf(`unable to upgrade %v`, err)
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)
	go client.ReadMessages()
	go client.WriteMessages()
}

func (m *Manager) addClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client] = true
	log.Printf(`clients %v`, m.clients)
}

func (m *Manager) removeClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
}