// // Inspired by https://github.com/gorilla/websocket/blob/master/examples/command/main.go
package websockets

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	wsconn         *websocket.Conn
	eventReceiver  chan *SocketEvent
	deregistration chan *Connection
	messageSender  chan []byte
	messageQueue   chan []byte
	close          chan bool
	timestamp      time.Time
	uuid           string
	connected      bool
}

const (
	maxMessageSize = 2048
	flushTime      = time.Millisecond * 30
	pingTime       = time.Second * 60
	pongTime       = pingTime + time.Second*10
	writeDelay     = time.Second * 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewConnection(conn *websocket.Conn, eventReceiver chan *SocketEvent, deregistration chan *Connection) *Connection {
	connection := &Connection{
		wsconn:         conn,
		eventReceiver:  eventReceiver,
		deregistration: deregistration,
		timestamp:      time.Now(),
		messageSender:  make(chan []byte, 256),
		messageQueue:   make(chan []byte, 256),
		close:          make(chan bool, 1),
		uuid:           uuid.New().String(),
	}
	return connection
}

func (connection *Connection) UUID() string {
	return connection.uuid
}

func (connection *Connection) Run() {
	connection.receiver()
	connection.sender()
	connection.connected = true

	<-connection.close
	connection.connected = false

	connection.deregistration <- connection
}

func (connection *Connection) Close() {
	if connection.connected {
		connection.close <- true
	}
	connection.wsconn.Close()
	defer close(connection.close)
	defer close(connection.messageSender)
	defer close(connection.messageQueue)
}

func (connection *Connection) Send(message []byte) {
	connection.messageSender <- message
}

func (connection *Connection) receiver() {
	connection.wsconn.SetReadLimit(maxMessageSize)
	connection.wsconn.SetReadDeadline(time.Now().Add(pongTime))
	connection.wsconn.SetPongHandler(func(string) error {
		connection.wsconn.SetReadDeadline(time.Now().Add(pongTime))
		return nil
	})
	go func() {
		for {
			_, data, err := connection.wsconn.ReadMessage()
			if err != nil {
				fmt.Println("Connection read error:", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					fmt.Println("Connection encountered an unexpected close error.") // These are the worst. Ugh.
				}
				break
			} else {
				data = bytes.TrimSpace(bytes.Replace(data, newline, space, -1))
				connection.eventReceiver <- &SocketEvent{
					Type:       SocketEventMessage,
					Connection: connection,
					Data:       data,
				}
			}
		}
		connection.close <- true
	}()
}

func (connection *Connection) sender() {
	go func() {
		flush := time.NewTicker(flushTime)
		defer flush.Stop()

		var buf bytes.Buffer
		for {
			select {
			case <-flush.C:
				message, err := buf.ReadBytes('\n')
				if err == nil {
					connection.messageQueue <- message
				}

			case message, ok := <-connection.messageSender:
				if !ok {
					fmt.Println("Connection send error:", message)
					return
				}
				buf.Write(message)
				buf.WriteString("\n")
			}
		}
	}()
	go func() {
		ping := time.NewTicker(pingTime)
		defer ping.Stop()

		for {
			select {
			case <-ping.C:
				connection.wsconn.SetWriteDeadline(time.Now().Add(writeDelay))
				if err := connection.wsconn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					fmt.Println("Connection write error:", err)
					return
				}

			case message, ok := <-connection.messageQueue:
				connection.wsconn.SetWriteDeadline(time.Now().Add(writeDelay))
				if !ok {
					connection.wsconn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := connection.wsconn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(message)
				if err := w.Close(); err != nil {
					fmt.Println("Connection write error:", err)
					return
				}
			}
		}
	}()
}
