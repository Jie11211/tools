package conf_test

import (
	"fmt"
	"testing"

	"github.com/Jie11211/tools/cicdtool/conf"
)

func TestLoadConfig(t *testing.T) {
	config := conf.LoadConfig("./config.yaml")
	fmt.Print(config)
}
