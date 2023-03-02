package tracetool

import (
	"context"
	"net/http"
)

type Trace interface {
	Start(ctx context.Context, name string) (context.Context, interface{})                                    // Start 开启第一个父span,interface{} 为span
	Transmit(ctx context.Context, span interface{}, name string) (context.Context, interface{})               // Transmit 函数中间传递 interface{} 为span
	Finish(span interface{})                                                                                  // Finish 结束span并上报
	SetSpanTag(span interface{}, tags map[string]string)                                                      // SetSpanTag 设置tag
	InjectHTTP(c context.Context, span interface{}, header *http.Header) error                                // InjectHTTP http传id
	ExtractHTTP(c context.Context, spanName string, header http.Header) (context.Context, interface{}, error) // ExtractHTTP 解析id interface{} 为span
	ExtractText(spanName, key, val string) (interface{}, error)                                               // ExtractText 解析id key:header的key val:header对应key的值 interface{} 为span
	RunHTTP(addr string, handler http.Handler) error                                                          // RunHTTP 启动http服务
	OnError(span interface{}, err error)                                                                      //设置错误
}
