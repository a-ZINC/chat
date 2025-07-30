package main

import (
	"encoding/json"
	"fmt"
	"log"

	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

var Clients map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager
	egress  chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan []byte),
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		messageType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf(`error connection closed %v`, err)
			}
			break
		}
		fmt.Printf(`msg %v %s`, messageType, msg)
		for val := range c.manager.clients {
			if val == c {
				continue
			}

			ms, err := json.Marshal(Message{
				Text:     string(msg),
				Username: "system",
			})
			if err != nil {
				continue
			}
			val.egress <- ms
		}
	}
}

func (c *Client) WriteMessages() {
	ticker := time.Ticker{}
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed.")
				}
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("failed to send message")
			}
		case <-ticker.C:
			fmt.Printf("hii")
		}
	}
}
