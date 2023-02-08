package httptool_test

import (
	"fmt"
	"testing"

	"github.com/Jie11211/tools/httptool"
)

func TestHttptool(t *testing.T) {
	h := httptool.NewHttptool()
	// i, err := h.DefaultPost("127.0.0.1:1236/test", httptool.ContentTypeJson, `{"test":1}`)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(i))
	h.AddParm("test", 111).AddParm("111", 111)
	h.AddHeader("xxx", "xxx")
	// s, err := h.Post("127.0.0.1:1236/test", httptool.ContentTypeJson)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(s))
	h.AddHeader("test", "aa")
	x, err := h.Get("127.0.0.1:1236/test", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(x))

}
