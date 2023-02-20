package client

import (
    "context"
    "fmt"
    v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {

	//1.client Set连接
	ClientSet *kubernetes.Clientset

	//2.dynamicClient

	DynamicClient *dynamic.DynamicClient

	//3.发现连接
	DiscoverryClient *discovery.DiscoveryClient

	//4.发现连接下的缓存
	DiskClient *disk.CachedDiscoveryClient
	//router  pre do post

}

func NewClient(ConfigPath string) *K8sClient {
    config, err := clientcmd.BuildConfigFromFlags("",ConfigPath)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	discoverryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	//    diskClient, err := disk.NewCachedDiscoveryClientForConfig(config, "","", time.Second)
	//    if err!=nil{
	//        panic(err)
	//    }
	return &K8sClient{
		ClientSet:        clientSet,
		DynamicClient:    dynamicClient,
		DiscoverryClient: discoverryClient,
		//        DiskClient: diskClient,
	}
}

func (k8sClient *K8sClient) A() {
	//创建deploy
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

    create, err := k8sClient.ClientSet.AppsV1().Deployments("default").Create(context.TODO(), deployManifest, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(create)
}
