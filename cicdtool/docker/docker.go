package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"k8s.io/apimachinery/pkg/util/json"
)

type DockerClient struct {
	Client *client.Client
}

func NewClient() *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &DockerClient{
		Client: cli,
	}
}

func (dc *DockerClient) Pull(ctx context.Context, refStr string, flag ...string) error {
	options := types.ImagePullOptions{}
	if len(flag) >= 2 {
		authConfig := types.AuthConfig{
			Username: flag[0],
			Password: flag[1],
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		options.RegistryAuth = authStr
	}
	out, err := dc.Client.ImagePull(ctx, refStr, options)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}

func (dc *DockerClient) Push(ctx context.Context, image string, flag ...string) (string, error) {
	options := types.ImagePushOptions{}
	if len(flag) >= 2 {
		authConfig := types.AuthConfig{
			Username: flag[0],
			Password: flag[1],
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return "", err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		options.RegistryAuth = authStr
	}
	reader, err := dc.Client.ImagePush(ctx, image, options)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	// 回显字符
	return buf.String(), nil
}

func (dc *DockerClient) CreateAndStart(ctx context.Context, config *container.Config) {
	resp, err := dc.Client.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := dc.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

// BuildWithDockerFile srcPath为Dockerfile的路径
func (dc *DockerClient) BuildWithDockerFile(srcPath string, tags []string) (string, error) {
	reader, _ := archive.TarWithOptions(srcPath, &archive.TarOptions{})
	resp, err := dc.Client.ImageBuild(context.Background(), reader, types.ImageBuildOptions{
		Dockerfile:  "Dockerfile",
		Tags:        tags,
		ForceRemove: true,
		NoCache:     true,
	})
	if err != nil {
		return "", err
	}
	//未读完程序结束的话不会创建
	buffer, err := io.ReadAll(resp.Body)
	return string(buffer), err
}

// 192.168.187.131:5000/alpinelinux/golang
func (dc *DockerClient) FindTag(ctx context.Context, repoTag string) (version string) {
	images, err := dc.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Println(image.RepoTags)
		//        if repoTag==image.RepoTags[0]{
		//            fmt.Println(image.RepoTags)
		//            version=strings.Split(image.RepoTags)
		//        }
	}
	return ""
}
