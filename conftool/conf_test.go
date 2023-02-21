package conftool_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Jie11211/tools/conftool"
)

type S struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var y = `{"name":"wyk","age":10}`

func TestConf(t *testing.T) {
	c := conftool.NewConf()
	m := c.LoadConfigWithStr(y, conftool.Json)
	var s S
	b, _ := json.Marshal(m)
	json.Unmarshal(b, &s)
	fmt.Println(s)
}
