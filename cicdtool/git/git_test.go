package git_test

import (
	"fmt"
	"testing"

	"github.com/Jie11211/tools/cicdtool/git"
)

func TestGitClient_Pull(t *testing.T) {
	client := git.NewClient("/root/.ssh/id_ed25519")
	_, err := client.Pull("git@gitee.com:xiayuyin/k8s-svc.git")
	if err != nil {
		fmt.Println(err)
	}
}
