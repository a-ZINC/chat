package main

import (
	"encoding/json"
	"fmt"
	"log"

	"time"

	"github.com/gorilla/websocket"
)

var Clients map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager
	egress  chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan Event),
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
		fmt.Printf("msg %v %s \n", messageType, msg)
		for val := range c.manager.clients {
			if val == c {
				continue
			}
			event := Event{}
			err := json.Unmarshal(msg, &event)
			if err != nil {
				fmt.Printf("error sending message %v", err)
				continue
			}
			val.egress <- event
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
					return
				}
				return
			}
			c.manager.eventRoute(message, c)
			
		case <-ticker.C:
			fmt.Printf("hii")
		}
	}
}