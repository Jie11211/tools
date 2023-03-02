package etcd

import (
	"context"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdCli struct {
	EtcdClient *clientv3.Client
}

func NewEtcdCli(Endpoints []string, Username, Password string, timeOut int) (*EtcdCli, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   Endpoints,
		DialTimeout: time.Duration(timeOut) * time.Second,
		Username:    Username,
		Password:    Password,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdCli{
		EtcdClient: etcdClient,
	}, nil
}

func (e *EtcdCli) Close() {
	e.EtcdClient.Close()
}

func (e *EtcdCli) Get(key string) string {
	value, err := e.EtcdClient.Get(context.Background(), key)
	if err != nil {
		return ""
	}
	for _, kv := range value.Kvs {
		return string(kv.Value)
	}
	return ""
}

func (e *EtcdCli) GetPrefixAll(Prefix string) (all map[string]string) {
	all = make(map[string]string)
	value, err := e.EtcdClient.Get(context.Background(), Prefix, clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, kv := range value.Kvs {
		all[string(kv.Key)] = string(kv.Value)
	}
	return
}

func (e *EtcdCli) Put(key string, value string) error {
	_, err := e.EtcdClient.Put(context.Background(), key, value)
	return err
}

func (e *EtcdCli) WatchPrefix(key string, prefix bool, Fput func(e *clientv3.Event) error, Fdel func(e *clientv3.Event) error) {
	var watchChan clientv3.WatchChan
	if prefix {
		watchChan = e.EtcdClient.Watch(context.TODO(), key, clientv3.WithPrefix())
	} else {
		watchChan = e.EtcdClient.Watch(context.TODO(), key, clientv3.WithKeysOnly())
	}
	for v := range watchChan {
		for _, event := range v.Events {
			switch event.Type {
			case mvccpb.PUT: //修改或者新增
				Fput(event)
			case mvccpb.DELETE: //删除
				Fdel(event)
			}
		}
	}
}
