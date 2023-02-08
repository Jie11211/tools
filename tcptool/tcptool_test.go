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

	c := tcptool.NewClient()
	err := c.Dial("127.0.0.1:9632")
	if err != nil {
		fmt.Println(err)
	}
	dp := tcptool.NewDataPack()
	msg := &tcptool.Msg{
		Id:     0,
		Data:   []byte("wyk"),
		MsgLen: uint32(len([]byte("wyk"))),
	}
	b := dp.Pack(*msg)

	for i := 0; i < 3; i++ {
		_, err = c.Conn.Write(b)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)
	}
	c.Conn.Close()
	fmt.Println(s.Conns)
}
