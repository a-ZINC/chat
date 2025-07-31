package main

import "encoding/json"

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

const (
	message = "message"
	system  = "system"
)

type EventHandler func(eve Event, c *Client) error
