package mysqltool

import (
	"fmt"
	"testing"
)

func TestConn(t *testing.T) {
	m := NewMysqlTool("root", "123456", "101.43.248.147", 3310, "end")
	err := m.Connect()
	if err != nil {
		fmt.Println(err)
	}
	defer m.Close()
	
}
