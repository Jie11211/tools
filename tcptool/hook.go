package tcptool

import "net"

type Hook interface {
	Pre(net.Conn, Msg) error
	Do(net.Conn, Msg) error
	Post(net.Conn) error
}

//特定场景上线广播
type DefaultHook struct {
	Ser *Server
}

func (dh *DefaultHook) Pre(_ net.Conn, msg Msg) error {
	dp := NewDataPack()
	b := dp.Pack(msg)
	for _, c := range dh.Ser.Conns {
		c.Conn.Write(b)
	}
	return nil
}

func (dh *DefaultHook) Do(net.Conn, Msg) error {
	return nil
}

func (dh *DefaultHook) Post(net.Conn) error {
	return nil
}
