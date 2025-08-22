package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Server struct {
	clients   map[*Client]bool
	clientsMu sync.RWMutex
	broadcast chan []byte
	upgrader  websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	s.clientsMu.Lock()
	s.clients[client] = true
	s.clientsMu.Unlock()

	go client.writePump()
	go client.readPump(s)
}

func (s *Server) BroadcastMessage(message interface{}) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	log.Printf("Broadcasting message to %d clients: %s", len(s.clients), string(msgBytes))

	s.broadcast <- msgBytes
}

func (s *Server) Run() {
	for {
		select {
		case message := <-s.broadcast:
			s.clientsMu.RLock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.clientsMu.RUnlock()
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(s *Server) {
	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, c)
		s.clientsMu.Unlock()
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
