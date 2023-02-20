package websocket

import (
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// host:127.0.0.0:2257  path: /test
func (wc *Client) NewConn(host, path string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
