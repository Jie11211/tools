package tracetool

import (
	"context"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

type JaegerTrace struct {
	tracer opentracing.Tracer
}

func (jt *JaegerTrace) Start(ctx context.Context, name string) (context.Context, interface{}) {
	span := jt.tracer.StartSpan(name)
	return ctx, span
}

func (jt *JaegerTrace) Transmit(ctx context.Context, span interface{}, name string) (context.Context, interface{}) {
	childSpan := jt.tracer.StartSpan(name, opentracing.ChildOf(span.(opentracing.Span).Context()))
	return ctx, childSpan
}

func (jt *JaegerTrace) Finish(span interface{}) {
	oSpan := span.(opentracing.Span)
	oSpan.Finish()
}

// SetSpanTag 设置tag
func (jt *JaegerTrace) SetSpanTag(span interface{}, tags map[string]string) {
	for k, v := range tags {
		span.(opentracing.Span).SetTag(k, v)
	}
}

// InjectHTTP http传id
func (jt *JaegerTrace) InjectHTTP(c context.Context, span interface{}, header *http.Header) error {
	return jt.tracer.Inject(span.(opentracing.Span).Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*header))
}

// ExtractHTTP 解析id
func (jt *JaegerTrace) ExtractHTTP(c context.Context, spanName string, header http.Header) (context.Context, interface{}, error) {
	spanContext, err := jt.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	childSpan := jt.tracer.StartSpan(spanName, opentracing.ChildOf(spanContext))
	return c, childSpan, err
}

// ExtractText 解析id key:header的key val:header对应key的值
func (jt *JaegerTrace) ExtractText(spanName, key, val string) (interface{}, error) {
	spanContext, err := jt.tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier{key: val})
	childSpan := jt.tracer.StartSpan(spanName, opentracing.ChildOf(spanContext))
	return childSpan, err
}

// RunHTTP 启动http服务
func (jt *JaegerTrace) RunHTTP(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, nethttp.Middleware(jt.tracer, handler))
}

// OnError 设置错误
func (jt *JaegerTrace) OnError(span interface{}, err error) {
	span.(opentracing.Span).SetTag(string(ext.Error), true)
	span.(opentracing.Span).LogKV(otlog.Error(err))
}
