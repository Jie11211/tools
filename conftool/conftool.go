package conftool

import (
	"strings"

	"github.com/spf13/viper"
)

type Conf struct {
}

func NewConf() *Conf {
	return &Conf{}
}

func (cf *Conf) LoadConfigWithPath(path string, confType ...confType) map[string]interface{} {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if len(confType) > 0 {
		viper.SetConfigType(string(confType[0]))
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	conf := make(map[string]interface{})
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return conf
}

func (cf *Conf) LoadConfigWithStr(data string, confType ...confType) map[string]interface{} {
	viper.SetConfigType("yaml")
	if len(confType) > 0 {
		viper.SetConfigType(string(confType[0]))
	}
	r := strings.NewReader(data)
	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}
	conf := make(map[string]interface{})
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return conf
}
