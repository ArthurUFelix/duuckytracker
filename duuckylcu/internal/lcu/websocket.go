package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn     *websocket.Conn
	password string
	baseURL  string
}

func NewWebsocketClient(port, password string) *WebSocketClient {
	baseURL := fmt.Sprintf("https://127.0.0.1:%s", port)

	return &WebSocketClient{
		baseURL:  baseURL,
		password: password,
	}
}

func (ws *WebSocketClient) Connect() error {
	wsURL := "wss" + ws.baseURL[5:] + "/"

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		HandshakeTimeout: 10 * time.Second,
	}

	auth := base64.StdEncoding.EncodeToString(([]byte("riot:" + ws.password)))
	header := http.Header{}
	header.Add("Authorization", "Basic "+auth)

	log.Printf("connection to LCU WebSocket: %s", wsURL)

	conn, _, err := dialer.Dial(wsURL, header)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	ws.conn = conn
	log.Println("WebSocket connected!")

	return nil
}

func (ws *WebSocketClient) Subscribe(eventName string) error {
	subscribeMsg := []interface{}{5, eventName}

	if err := ws.conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", eventName, err)
	}

	log.Printf("subscribed to event: %s", eventName)
	return nil
}

func (ws *WebSocketClient) Listen(callback func(LCUEvent)) error {
	for {
		_, messageBytes, err := ws.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		if len(messageBytes) == 0 {
			continue
		}

		var rawMsg []interface{}
		if err := json.Unmarshal(messageBytes, &rawMsg); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		if len(rawMsg) < 3 {
			continue
		}

		msgType, ok := rawMsg[0].(float64)
		if !ok || msgType != 8 {
			continue
		}

		uri, ok := rawMsg[1].(string)
		if !ok {
			continue
		}

		eventData, ok := rawMsg[2].(map[string]interface{})
		if !ok {
			continue
		}

		dataBytes, err := json.Marshal(eventData)
		if err != nil {
			log.Printf("failed to marshal event data: %v", err)
			continue
		}

		var event LCUEvent
		if err := json.Unmarshal(dataBytes, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			continue
		}

		event.URI = uri
		callback(event)
	}
}

func (ws *WebSocketClient) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}
