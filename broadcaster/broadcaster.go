package broadcaster

import (
	"github.com/googollee/go-socket.io"
	"log"
)

type socketIoServer struct {
	Server *socketio.Server
}

func (h socketIoServer) Broadcast(message string) {
	log.Println("broadcasting", message)
	h.Server.BroadcastTo("climate", message)
}

func New() socketIoServer {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatalln(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("Connection")
		so.Join("climate")
		so.On("disconnection", func() {
			log.Println("disconnected")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("Error", err)
	})
	return socketIoServer{Server: server}
}
