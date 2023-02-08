package tcptool

const (
	HeadLen = 8 // uint32 占四个字节 Id+MsgLen 共占8个字节
)

type Msg struct {
	Id     uint32 //类型，对应处理方法
	Data   []byte
	MsgLen uint32
}

func NewInitMsg(data string) *Msg {
	data = data + "，已连接"
	b := []byte(data)
	return &Msg{
		Id:     0,
		Data:   b,
		MsgLen: uint32(len(b)),
	}
}

func NewMsg(id uint32, data string) *Msg {
	b := []byte(data)
	return &Msg{
		Id:     id,
		Data:   b,
		MsgLen: uint32(len(b)),
	}
}
