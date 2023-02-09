package tcptool

import (
	"fmt"
	"net"
)

type Client struct {
	Hook map[uint32]Hook
	Conn net.Conn
}

func NewClient() *Client {
	return &Client{
		Hook: make(map[uint32]Hook),
	}
}

func (c *Client) Dial(hostPort string) error {
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Client) Write(dp DataPack, id uint32, data string) error {
	m := NewMsg(id, data)
	b := dp.Pack(*m)
	_, err := c.Conn.Write(b)
	return err
}

func (c *Client) Read(dp DataPack, conn net.Conn) {
	for {
		m, err := dp.UnPackWithConn(conn, maxLen)
		if err != nil {
			fmt.Println(err)
			if err == ErrConnField {
				conn.Close()
				break
			}
		}
		//处理逻辑
		if hook, ok := c.Hook[m.Id]; ok {
			hook.Do(conn, *m)
		}
	}
}

func (c *Client) AddHook(id uint32, hook Hook) error {
	if _, ok := c.Hook[id]; !ok {
		c.Hook[id] = hook
		return nil
	}
	return ErrHookExist
}
