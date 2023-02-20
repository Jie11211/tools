package tools_test

import (
	"fmt"
	"testing"

	"github.com/Jie11211/tools/cicdtool/tools"
)

func TestGetNewVersion(t *testing.T) {
	version := tools.GetNewVersion("v0.0.0")
	fmt.Println(version)
}
