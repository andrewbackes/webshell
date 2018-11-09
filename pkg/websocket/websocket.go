package websocket

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"net/http"
	"sync"
	"time"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Server used to manage communication over websockets.
type Server struct {
	sync.RWMutex
	websockets map[string]*websocket.Conn
	upgrader   websocket.Upgrader
	msgHandler func(Message)
}

// Message received from a websocket.
type Message struct {
	ID        string
	Timestamp time.Time
	Value     []byte
	From      string
}

func (m Message) String() string {
	return fmt.Sprintf("{Timestamp:%v, ID:%v, From:%v, Value:%v}", m.Timestamp, m.ID, m.From, string(m.Value))
}

// Option for constructing a Server.
type Option func(*Server) error

// MessageHandler sets the function to be called when a message is received.
func MessageHandler(h func(Message)) Option {
	return func(s *Server) error {
		s.msgHandler = h
		return nil
	}
}

// NewServer for managing websockets.
func NewServer(opts ...Option) *Server {
	defaultMessageHandler := func(m Message) {
		fmt.Println(m)
	}
	s := &Server{
		msgHandler: defaultMessageHandler,
		websockets: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// UpgradeHandler upgrades a connection to a websocket.
func (s *Server) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := uuid.New()
	fmt.Printf("Client %v connected.\n", id)
	s.Lock()
	s.websockets[id] = conn
	s.Unlock()
	go func() {
		exists := func() bool {
			_, exists := s.websockets[id]
			return exists
		}
		for exists() {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("Websocket read error: %v", err)
				}
				break
			}
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			s.msgHandler(Message{
				Timestamp: time.Now(),
				ID:        uuid.New(),
				Value:     message,
				From:      id,
			})
		}
		fmt.Printf("Stopped listening to client %v.\n", id)
		s.terminate(id)
	}()

}

// Write a message to all websockets.
func (s *Server) Write(data []byte) {
	s.RLock()
	for id, conn := range s.websockets {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Printf("Error - could not write to websocket. Closing %v.\n", id)
			s.terminate(id)
		}
	}
	s.RUnlock()
}

func (s *Server) terminate(connID string) {
	s.websockets[connID].Close()
	s.Lock()
	delete(s.websockets, connID)
	s.Unlock()
	fmt.Printf("Terminated connection to %v.\n", connID)
}
