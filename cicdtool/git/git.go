package git

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitClient struct {
	PemFile string
}

// git@gitee.com:xiayuyin/k8s-svc.git
func NewClient(PemFile string) *GitClient {
	return &GitClient{
		PemFile: PemFile,
	}
}

// os.RemoveAll(temp) 处理完逻辑后，删除文件夹
func (gc *GitClient) Pull(url string) (string, error) {
	split := strings.Split(url, ".")
	names := strings.Split(split[len(split)-2], "/")
	temp, err := os.MkdirTemp("", fmt.Sprintf("%s-%d", names[len(names)-1], time.Now().Unix()))
	if err != nil {
		fmt.Println(err)
	}
	pri, err := ssh.NewPublicKeysFromFile("git", gc.PemFile, "") //privateFilePath为文件路径,内容如上
	if err != nil {
		return "", err
	}
	_, err = git.PlainClone(temp, false, &git.CloneOptions{
		URL:  url,
		Auth: pri,
	})
	return temp, nil
}
