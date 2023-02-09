package tcptool_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Jie11211/tools/tcptool"
)

type testHook struct{}

func (dh *testHook) Pre(_ net.Conn, msg tcptool.Msg) error {
	return nil

}

func (dh *testHook) Do(_ net.Conn, msg tcptool.Msg) error {
	fmt.Println(string(msg.Data))
	return nil
}

func (dh *testHook) Post(_ net.Conn) error {
	return nil
}

func TestPack(t *testing.T) {
	s := tcptool.NewServer(false)
	s.AddHook(0, &testHook{})
	go s.ListenAndAccept(":9632")
	time.Sleep(time.Second)
	dp := tcptool.NewDataPack()
	msg := &tcptool.Msg{
		Id:     0,
		Data:   []byte("server wyk"),
		MsgLen: uint32(len([]byte("server wyk"))),
	}
	b := dp.Pack(*msg)
	go func(ts *tcptool.Server, bs []byte) {
		for {
			for _, connTool := range ts.Conns {
				connTool.Conn.Write(bs)
			}
		}
	}(s, b)
	c := tcptool.NewClient()
	c.AddHook(0, &testHook{})
	err := c.Dial("127.0.0.1:9632")
	if err != nil {
		fmt.Println(err)
	}
	msg = &tcptool.Msg{
		Id:     0,
		Data:   []byte("wyk"),
		MsgLen: uint32(len([]byte("wyk"))),
	}
	b = dp.Pack(*msg)
	go func(dp tcptool.DataPack, tc *tcptool.Client) {
		c.Read(dp, c.Conn)
	}(*dp, c)

	for {
		_, err = c.Conn.Write(b)
		if err != nil {
			fmt.Println(err)
		}
	}
}
