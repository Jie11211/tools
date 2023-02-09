package tcptool

type TcpTool struct {
	Server *Server
	Client *Client
}

func NewTcpTool() *TcpTool {
	return &TcpTool{
		Server: NewServer(false),
		Client: NewClient(),
	}
}
