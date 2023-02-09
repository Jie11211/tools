package filetool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Jie11211/tools/filetool"
)

func TestFile(t *testing.T) {
	f, err := filetool.OpenFile(``)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	s, err := filetool.ReadLineString(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(time.Now())
	for _, v := range s {
		fmt.Println(v)
	}
	fmt.Println(time.Now())
}
