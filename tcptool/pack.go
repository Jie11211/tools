package tcptool

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

var (
	ErrMsgToLong = errors.New("MsgToLong")
	ErrConnField = errors.New("MsgConnField")
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) Pack(msg Msg) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, msg.MsgLen)
	binary.Write(bytesBuffer, binary.BigEndian, msg.Id)
	binary.Write(bytesBuffer, binary.BigEndian, msg.Data)
	return bytesBuffer.Bytes()
}

func (dp *DataPack) unPackHead(msgBtye []byte) Msg {
	var msg Msg
	r := bytes.NewReader(msgBtye)
	binary.Read(r, binary.BigEndian, &msg.MsgLen)
	binary.Read(r, binary.BigEndian, &msg.Id)
	return msg
}

func (dp *DataPack) UnPackWithConn(conn net.Conn, maxLen uint32) (*Msg, error) {
	head := make([]byte, HeadLen)
	_, err := io.ReadFull(conn, head)
	if err != nil {
		return nil, ErrConnField
	}
	m := dp.unPackHead(head)
	if m.MsgLen > 0 && maxLen >= m.MsgLen {
		data := make([]byte, m.MsgLen)
		_, err := io.ReadFull(conn, data)
		if err != nil {
			return nil, err
		}
		m.Data = data
	} else {
		return nil, ErrMsgToLong
	}
	return &m, nil
}
