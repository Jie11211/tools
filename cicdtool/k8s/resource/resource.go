package  resource

import (
    "context"
    "github.com/Jie11211/tools/cicdtool/k8s/client"
    "fmt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/apimachinery/pkg/util/json"
)

type Resource struct {
    Name string
    Kind string
    NameSpace string
    RsConfig interface{}
    //GVR
}

func (rs *Resource)Create(client *client.K8sClient)error{
    res, err :=runtime.DefaultUnstructuredConverter.ToUnstructured(rs.RsConfig)
    if err!=nil{
        fmt.Println(err)
    }
    obj := &unstructured.Unstructured{Object: res}
    _, APIResourceList, err := client.DiscoverryClient.ServerGroupsAndResources()
    //遍历
    //自定义好kind eg: pod
    groupVersionResource:=&schema.GroupVersionResource{}
    for gvIdx := 0; gvIdx < len(APIResourceList); gvIdx++ {
        //把 groupVersion切割开  gv的用json.Marshal出来的值{"Group":"networking.k8s.io","Version":"v1"} ；APIResourceList[gvIdx]的groupVersion值为 "groupVersion":"networking.k8s.io/v1"
        gv, err := schema.ParseGroupVersion(APIResourceList[gvIdx].GroupVersion)
        if err != nil {
            return err
        }
        //获取具体的资源
        for rIdx := 0; rIdx < len(APIResourceList[gvIdx].APIResources); rIdx++ {
            if rs.Kind == APIResourceList[gvIdx].APIResources[rIdx].Kind {
                groupVersionResource.Group=gv.Group
                groupVersionResource.Version=gv.Version
                groupVersionResource.Resource=APIResourceList[gvIdx].APIResources[rIdx].Name
                break
            }
        }
    }
    //K8sCliDynamic  DynamicClient客户端
    ax, err := client.DynamicClient.Resource(*groupVersionResource).Namespace(obj.GetNamespace()).Get(context.TODO(), rs.Name, metav1.GetOptions{})
    if err!=nil{
       fmt.Println(err)
       info, err := client.DynamicClient.Resource(*groupVersionResource).Namespace(obj.GetNamespace()).Create(context.TODO(), obj, metav1.CreateOptions{})
       if info.Object!=nil{
       fmt.Println("success")
       }
       return err
    }
    marshal, err := json.Marshal(ax.Object)
    fmt.Println(string(marshal))
    _, err = client.DynamicClient.Resource(*groupVersionResource).Namespace(obj.GetNamespace()).Update(context.TODO(), obj, metav1.UpdateOptions{})
    return err
}