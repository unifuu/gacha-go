package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"gacha/models"
	"gacha/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Message types
const (
	TypeSinglePull   = "single_pull"
	TypeTenPull      = "ten_pull"
	TypeGetUserInfo  = "get_user_info"
	TypeGetInventory = "get_inventory"
	TypeGetPool      = "get_pool"
	TypeAddCurrency  = "add_currency"

	// Response types
	TypeGachaResult    = "gacha_result"
	TypeUserInfo       = "user_info"
	TypeInventory      = "inventory"
	TypePoolInfo       = "pool_info"
	TypeCurrencyUpdate = "currency_update"
	TypeError          = "error"
	TypePing           = "ping"
	TypePong           = "pong"
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	gachaService *services.GachaService
	userService  *services.UserService
	clients      map[*websocket.Conn]*Client
	clientsMu    sync.RWMutex
}

// Client represents a connected WebSocket client
type Client struct {
	conn     *websocket.Conn
	username string
	send     chan []byte
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(gachaService *services.GachaService, userService *services.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		gachaService: gachaService,
		userService:  userService,
		clients:      make(map[*websocket.Conn]*Client),
	}
}

// HandleWebSocket handles WebSocket connection
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		conn:     conn,
		username: "default", // Can be extended to get from query params
		send:     make(chan []byte, 256),
	}

	h.clientsMu.Lock()
	h.clients[conn] = client
	h.clientsMu.Unlock()

	log.Printf("New WebSocket client connected: %s", conn.RemoteAddr())

	// Start goroutines for reading and writing
	go h.receiveMessages(client)
	go h.writeMessages(client)

	// Send initial user info
	h.sendUserInfo(client)
}

// receiveMessages handles incoming messages from the client
func (h *WebSocketHandler) receiveMessages(client *Client) {
	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, client.conn)
		h.clientsMu.Unlock()
		client.conn.Close()
		log.Printf("Client disconnected: %s", client.conn.RemoteAddr())
	}()

	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		h.handleMessage(client, message)
	}
}

// writeMessages handles outgoing messages to the client
func (h *WebSocketHandler) writeMessages(client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages
func (h *WebSocketHandler) handleMessage(client *Client, message []byte) {
	var msg WebSocketMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		h.sendError(client, "Invalid message format")
		return
	}

	switch msg.Type {
	case TypePing:
		h.sendPong(client)

	case TypeSinglePull:
		h.handleSinglePull(client)

	case TypeTenPull:
		h.handleTenPull(client)

	case TypeGetUserInfo:
		h.sendUserInfo(client)

	case TypeGetInventory:
		h.sendInventory(client)

	case TypeGetPool:
		h.sendPoolInfo(client)

	case TypeAddCurrency:
		var req models.AddCurrencyRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			h.sendError(client, "Invalid currency request")
			return
		}
		h.handleAddCurrency(client, req.Amount)

	default:
		h.sendError(client, "Unknown message type")
	}
}

// handleSinglePull processes single pull request
func (h *WebSocketHandler) handleSinglePull(client *Client) {
	user := h.userService.GetUser(client.username)
	if user == nil {
		h.sendError(client, "User not found")
		return
	}

	if !user.DeductCurrency(160) {
		h.sendError(client, "Insufficient currency")
		return
	}

	char := h.gachaService.PerformSinglePull(user)
	isNew := user.AddCharacter(char)

	result := models.GachaResult{
		Characters: []models.Character{char},
		IsNew:      []bool{isNew},
		Timestamp:  time.Now().Unix(),
	}

	h.sendMessage(client, TypeGachaResult, result)
	h.sendUserInfo(client)
}

// handleTenPull processes ten pull request
func (h *WebSocketHandler) handleTenPull(client *Client) {
	user := h.userService.GetUser(client.username)
	if user == nil {
		h.sendError(client, "User not found")
		return
	}

	if !user.DeductCurrency(1600) {
		h.sendError(client, "Insufficient currency")
		return
	}

	characters := h.gachaService.PerformTenPull(user)
	var isNewList []bool

	for _, char := range characters {
		isNew := user.AddCharacter(char)
		isNewList = append(isNewList, isNew)
	}

	result := models.GachaResult{
		Characters: characters,
		IsNew:      isNewList,
		Timestamp:  time.Now().Unix(),
	}

	h.sendMessage(client, TypeGachaResult, result)
	h.sendUserInfo(client)
}

// handleAddCurrency adds currency to user
func (h *WebSocketHandler) handleAddCurrency(client *Client, amount int) {
	user := h.userService.GetUser(client.username)
	if user == nil {
		h.sendError(client, "User not found")
		return
	}

	user.AddCurrency(amount)

	response := models.CurrencyResponse{
		Currency: user.Currency,
	}

	h.sendMessage(client, TypeCurrencyUpdate, response)
	h.sendUserInfo(client)
}

// sendUserInfo sends user information to client
func (h *WebSocketHandler) sendUserInfo(client *Client) {
	user := h.userService.GetUser(client.username)
	if user == nil {
		h.sendError(client, "User not found")
		return
	}

	response := models.UserInfoResponse{
		Username:  user.Username,
		Currency:  user.Currency,
		PityCount: user.PityCount,
	}

	h.sendMessage(client, TypeUserInfo, response)
}

// sendInventory sends user inventory to client
func (h *WebSocketHandler) sendInventory(client *Client) {
	user := h.userService.GetUser(client.username)
	if user == nil {
		h.sendError(client, "User not found")
		return
	}

	response := models.InventoryResponse{
		Inventory: user.Inventory,
		Count:     len(user.Inventory),
	}

	h.sendMessage(client, TypeInventory, response)
}

// sendPoolInfo sends pool information to client
func (h *WebSocketHandler) sendPoolInfo(client *Client) {
	poolInfo := models.PoolInfo{
		Characters: h.gachaService.GetCharacterPool(),
		Rates: map[string]string{
			"ssr": "2%",
			"sr":  "10%",
			"r":   "88%",
		},
		PitySystem: "Guaranteed SSR at 90 pulls",
	}

	h.sendMessage(client, TypePoolInfo, poolInfo)
}

// sendMessage sends a typed message to client
func (h *WebSocketHandler) sendMessage(client *Client, msgType string, data interface{}) {
	msg := WebSocketMessage{
		Type: msgType,
	}

	// Marshal data to JSON bytes
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marshal data: %v", err)
			return
		}
		msg.Data = jsonData
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	select {
	case client.send <- msgBytes:
	default:
		log.Printf("Client send channel full, dropping message")
	}
}

// sendError sends an error message to client
func (h *WebSocketHandler) sendError(client *Client, errMsg string) {
	msg := WebSocketMessage{
		Type:  TypeError,
		Error: errMsg,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal error: %v", err)
		return
	}

	select {
	case client.send <- msgBytes:
	default:
		log.Printf("Client send channel full, dropping error message")
	}
}

// sendPong sends a pong response
func (h *WebSocketHandler) sendPong(client *Client) {
	msg := WebSocketMessage{
		Type: TypePong,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal pong: %v", err)
		return
	}

	select {
	case client.send <- msgBytes:
	default:
		// Pong is not critical, can be dropped
	}
}
