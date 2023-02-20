package conf

import (
    "github.com/spf13/viper"
)
type Config struct {
    //git congif
   Git GitConfig `mapstructure:"GIT"`

    //docker config
    Docker DockerConfig`mapstructure:"Docker"`

    //k8s连接config
    K8s K8sConfig`mapstructure:"K8s"`


    //k8s资源config 切片
    Resource []ResourceConfig`mapstructure:"Resource"`
}

type GitConfig struct {
    Used bool`mapstructure:"Used"`
    PrivateFilePath string`mapstructure:"PrivateFilePath"`
}

type DockerConfig struct {
    Used bool`mapstructure:"Used"`
    ImageUrl string`mapstructure:"ImageUrl"`
    UserName string`mapstructure:"UserName"`
    PassWord string`mapstructure:"PassWord"`
}

type K8sConfig struct {
    ConfigPath string`mapstructure:"ConfigPath"`
}

type ResourceConfig struct {
    Kind string`mapstructure:"kind"`
    Replicase int`mapstructure:"replicase"`
}


func LoadConfig(path string)*Config{
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")
    if err := viper.ReadInConfig();err!=nil{
        panic(err)
    }
    conf:=&Config{}
    if err := viper.Unmarshal(conf);err!=nil{
        panic(err)
    }
    return conf
}
