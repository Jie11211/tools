package tcptool

import (
	"errors"
	"fmt"
	"net"
)

const (
	maxLen  = 1024
	maxConn = 100
)

var (
	ErrHookExist    = errors.New("HookExist")
	ErrConnNotExist = errors.New("ConnNotExist")
)

type Server struct {
	Conns map[string]ConnTool
	Hook  map[uint32]Hook
	Close map[net.Conn]string
}
type ConnTool struct {
	Conn net.Conn
	IP   string
}

func NewServer(defaulthook bool) *Server {
	s := &Server{
		Hook:  make(map[uint32]Hook),
		Conns: make(map[string]ConnTool),
		Close: make(map[net.Conn]string),
	}
	if defaulthook {
		s.Hook = map[uint32]Hook{
			1: &DefaultHook{},
		}
	}
	return s
}

// Accept hostPort可以为ip:port 也可以为 :port
func (s *Server) ListenAndAccept(hostPort string) {
	l, err := net.Listen("tcp", hostPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		if len(s.Conns) >= maxConn {
			continue
		}
		dp := NewDataPack()
		initConn(s, conn, dp)
		if conn == nil {
			continue
		}
		go s.ReadMsg(conn, dp)
	}
}

// 客户端第一次请求为初始化连接，建立 name：conn的对应关系
func initConn(s *Server, conn net.Conn, dp *DataPack) {
	for {
		m, err := dp.UnPackWithConn(conn, maxLen)
		if m.MsgLen == 0 || err != nil {
			if err == ErrConnField {
				conn.Close()
				break
			}
			continue
		}
		s.Conns[string(m.Data)] = ConnTool{Conn: conn, IP: conn.RemoteAddr().String()}
		s.Close[conn] = string(m.Data)
		if dh, ok := s.Hook[1].(*DefaultHook); ok {
			dh.Ser = s
			dh.Pre(conn, *NewInitMsg(string(m.Data)))
		}
		break
	}
}

func (s *Server) ReadMsg(conn net.Conn, dp *DataPack) {
	for {
		m, err := dp.UnPackWithConn(conn, maxLen)
		if err != nil {
			if err == ErrConnField {
				name := s.Close[conn]
				delete(s.Close, conn)
				delete(s.Conns, name)
				conn.Close()
				break
			}
		}
		//处理逻辑
		if hook, ok := s.Hook[m.Id]; ok {
			hook.Do(conn, *m)
		}
	}
}

func (s *Server) AddHook(id uint32, hook Hook) error {
	if _, ok := s.Hook[id]; !ok {
		s.Hook[id] = hook
		return nil
	}
	return ErrHookExist
}

func (s *Server) Write(dp DataPack, id uint, data string, name string) error {
	if connTool, ok := s.Conns[name]; ok {
		_, err := connTool.Conn.Write(dp.Pack(Msg{
			Id:     uint32(id),
			Data:   []byte(data),
			MsgLen: uint32(len([]byte(data))),
		}))
		return err
	}
	return ErrConnNotExist
}

// TODO 读写分离
