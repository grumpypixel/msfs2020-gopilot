// Inspired by https://github.com/gorilla/websocket/blob/master/examples/command/main.go
package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type SocketEvent struct {
	Type       int
	Connection *Connection
	Data       []byte
}

type WebSocket struct {
	EventReceiver  chan *SocketEvent
	connections    map[*Connection]int
	registration   chan *Connection
	deregistration chan *Connection
	broadcaster    chan []byte
	sender         chan *SocketEvent
}

const (
	SocketEventConnected = iota
	SocketEventDisconnected
	SocketEventMessage
)

var (
	upgrader websocket.Upgrader
)

func init() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

func NewWebSocket() *WebSocket {
	socket := &WebSocket{
		EventReceiver:  make(chan *SocketEvent),
		connections:    make(map[*Connection]int),
		registration:   make(chan *Connection),
		deregistration: make(chan *Connection),
		broadcaster:    make(chan []byte, 256),
		sender:         make(chan *SocketEvent),
	}
	go socket.Run()
	return socket
}

func (socket *WebSocket) ConnectionCount() int {
	return len(socket.connections)
}

func (socket *WebSocket) ConnectionUUIDs() []string {
	uuids := make([]string, 0)
	for key, _ := range socket.connections {
		uuids = append(uuids, key.UUID())
	}
	return uuids
}

func (socket *WebSocket) Run() {
	for {
		select {
		case connection := <-socket.registration:
			socket.connections[connection] = 1
			socket.EventReceiver <- &SocketEvent{
				Type:       SocketEventConnected,
				Connection: connection,
				Data:       nil,
			}

		case connection := <-socket.deregistration:
			if _, exists := socket.connections[connection]; exists {
				socket.EventReceiver <- &SocketEvent{
					Type:       SocketEventDisconnected,
					Connection: connection,
					Data:       nil,
				}
				connection.Close()
				delete(socket.connections, connection)
			}

		case data := <-socket.broadcaster:
			for client := range socket.connections {
				client.Send(data)
			}

		case message := <-socket.sender:
			message.Connection.Send(message.Data)
		}
	}
}

func (socket *WebSocket) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	connection := NewConnection(conn, socket.EventReceiver, socket.deregistration)
	socket.registration <- connection
	connection.Run()
}

func (socket *WebSocket) Broadcast(message []byte) {
	socket.broadcaster <- message
}

func (socket *WebSocket) Send(connectionUUID string, data []byte) bool {
	for connection := range socket.connections {
		if connection.uuid == connectionUUID {
			socket.sender <- &SocketEvent{
				Type:       SocketEventMessage,
				Connection: connection,
				Data:       data,
			}
			return true
		}
	}
	return false
}
