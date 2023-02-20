package client_test

import (
	"testing"

	"github.com/Jie11211/tools/cicdtool/k8s/client"
)

func TestNewClient(t *testing.T) {
	client.NewClient("/etc/kubernetes/admin.conf").A()
}
