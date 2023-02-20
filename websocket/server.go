package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	Upgrader *websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		Upgrader: &websocket.Upgrader{},
	}
}

func (ws *Server) Upgrad(w http.ResponseWriter, r *http.Request, header http.Header) (*websocket.Conn, error) {
	c, err := ws.Upgrader.Upgrade(w, r, header)
	if err != nil {
		return nil, err
	}
	return c, nil
}
