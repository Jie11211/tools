package resource_test

import (
	"fmt"
	"testing"

	"github.com/Jie11211/tools/cicdtool/k8s/client"
	"github.com/Jie11211/tools/cicdtool/k8s/resource"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResource_Create(t *testing.T) {
    k8sClient := client.NewClient("/etc/kubernetes/admin.conf")
    rs:=&resource.Resource{
        Name:      "nginx-test",
        Kind:      "Deployment",
        NameSpace: "default",
        RsConfig:  a(),
    }
    err := rs.Create(k8sClient)
    if err!=nil{
        fmt.Println(err)
    }
}

func a()*v1.Deployment{
    deployManifest := new(v1.Deployment) //配置文件
    deployManifest.Name = "nginx-test"
    deployManifest.Kind = "Deployment"
    deployManifest.Namespace = "default"
    deployManifest.APIVersion = "apps/v1"
    replicase := int32(1)
    deployManifest.Spec.Replicas = &replicase
    deployManifest.Spec.Selector = &metav1.LabelSelector{
        MatchLabels: map[string]string{
            "app": "nginx-test",
            },
            }
            deployManifest.Spec.Template.Labels = map[string]string{
                "app": "nginx-test",
                }
                deployManifest.Spec.Strategy = v1.DeploymentStrategy{
                    Type: v1.RollingUpdateDeploymentStrategyType,
                    }
                    deployManifest.Spec.Template.ObjectMeta.Labels = map[string]string{
                        "app": "nginx-test",
                        }

                        deployManifest.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways
                        c := corev1.Container{
                            Name:            "nginx-test",
                            Image:           "nginx",
                            ImagePullPolicy: corev1.PullIfNotPresent,
                            }
                            deployManifest.Spec.Template.Spec.Containers = append(deployManifest.Spec.Template.Spec.Containers, c)
                            deployManifest.Spec.Template.Spec.DNSPolicy = corev1.DNSClusterFirst
                            return deployManifest
}