package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 開発用：本番ではオリジンチェックを実装
	},
}

type Client struct {
	conn   *websocket.Conn
	roomID string
	userID string
	send   chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *types.Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *types.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run(msgStore store.IMessageStore) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: roomID=%s, userID=%s", client.roomID, client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: roomID=%s, userID=%s", client.roomID, client.userID)

		case message := <-h.broadcast:
			// メッセージを保存
			log.Printf("Saving message - RoomID: %s, UserID: %s, Content: %s",
				message.RoomID, message.UserID, message.Content)
			savedMsg, err := msgStore.SaveMessage(message.RoomID, message.UserID, message.Content)
			if err != nil {
				log.Printf("Failed to save message: %v", err)
				continue
			}

			log.Printf("Message saved, broadcasting to room %s. Active clients: %d",
				message.RoomID, len(h.clients))

			// 同じ部屋のクライアントにブロードキャスト
			messageData, _ := json.Marshal(savedMsg)
			h.mu.RLock()
			broadcastCount := 0
			for client := range h.clients {
				if client.roomID == message.RoomID {
					broadcastCount++
					select {
					case client.send <- messageData:
					default:
						log.Printf("Failed to send to client, closing connection")
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
			log.Printf("Broadcasted to %d clients in room %s", broadcastCount, message.RoomID)
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("ReadMessage error: %v", err)
			break
		}

		log.Printf("Received message: %s", string(message))

		var req types.SendMessageRequest
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("Invalid message format: %v, raw: %s", err, string(message))
			continue
		}

		log.Printf("Parsed request - Action: %s, RoomID: %s, UserID: %s, Content: %s",
			req.Action, req.RoomID, c.userID, req.Content)

		if req.Action == "sendMessage" {
			log.Printf("Broadcasting message to room %s", req.RoomID)
			hub.broadcast <- &types.Message{
				RoomID:  req.RoomID,
				UserID:  c.userID,
				Content: req.Content,
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	userID := r.URL.Query().Get("user_id")

	if roomID == "" || userID == "" {
		http.Error(w, "room_id and user_id are required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn:   conn,
		roomID: roomID,
		userID: userID,
		send:   make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

func main() {
	ctx := context.Background()
	msgStore := store.NewMessageStore(ctx, "ChatAppMessages")

	hub := newHub()
	go hub.run(msgStore)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("WebSocket server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
