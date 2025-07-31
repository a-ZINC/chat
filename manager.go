package main

import (
	"encoding/json"
	"fmt"
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
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients: make(map[*Client]bool),
		mu: sync.Mutex{},
		handlers: make(map[string]EventHandler),
	}
	m.setupHandler()
	return m
}

func(m *Manager) setupHandler() {
	m.handlers[message] = m.sendMessage
}

func(m *Manager) sendMessage(eve Event, client *Client) error {
	msg, err := json.Marshal(eve)
	if err != nil {
		return err
	}
	fmt.Printf("msg: %v", msg)
	return client.conn.WriteMessage(websocket.TextMessage, msg)
}

func(m *Manager) eventRoute(eve Event, client *Client) {
	fmt.Printf("eve %v", eve)
	if handler, ok := m.handlers[eve.Type]; ok {
		err := handler(eve, client)
		if err != nil {
			fmt.Printf("error sending Event through handler")
			return;
		}
	} else {
		fmt.Printf("error sending Event")
		return;
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
}

func (m *Manager) removeClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
	}
}