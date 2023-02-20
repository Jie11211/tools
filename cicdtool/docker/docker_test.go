package docker_test

import (
	"context"

	//    "context"
	"github.com/Jie11211/tools/cicdtool/docker"
	// "fmt"
	"testing"
	// "time"
)

func TestDockerClient_FindTag(t *testing.T) {
	docker.NewClient().FindTag(context.Background(), "")
}

//func TestDockerClient_BuildWithDockerFile(t *testing.T) {
//    client:= docker.NewClient()
//    file, err := client.BuildWithDockerFile("./", []string{"192.168.187.131:5000/alpinelinux/golang:v0.0.2"})
//    if err!=nil{
//        fmt.Println(err)
//    }
//    fmt.Println(file)
//}

//func TestNewClient(t *testing.T) {
//
//}

//func TestDockerClient_Pull(t *testing.T) {
//    client := docker.NewClient()
//    ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
//    defer cancelFunc()
//    err:= client.Pull(ctx, "192.168.187.131:5000/alpine:v0.0.1")
//    if err!=nil{
//        fmt.Println(err)
//    }
//}
//
//func TestDockerClient_Push(t *testing.T) {
//    client := docker.NewClient()
//    ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
//    defer cancelFunc()
//    _, err := client.Push(ctx, "192.168.187.131:5000/nginx:v0.0.1","","")
//    if err!=nil{
//        fmt.Println(err)
//    }
//}
