package tracetool

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultTraceId = "Trace-Id"
	defaultSpanId  = "Span-Id"
)

var (
	traceId = defaultTraceId
	spanId  = defaultSpanId
)

type OpenTelemetryTrace struct {
}

// Start 开启第一个父span
func (ot *OpenTelemetryTrace) Start(ctx context.Context, name string) (context.Context, interface{}) {
	tracer := otel.Tracer(name)
	spanCtx, span := tracer.Start(ctx, name)
	return spanCtx, span
}

// Transmit 函数中间传递
func (ot *OpenTelemetryTrace) Transmit(ctx context.Context, span interface{}, Name string) (context.Context, interface{}) {
	spanCtx, oSpan := otel.Tracer(Name).Start(ctx, Name)
	return spanCtx, oSpan
}

// Finish 结束span并上报
func (ot *OpenTelemetryTrace) Finish(span interface{}) {
	span.(trace.Span).End()
}

// SetSpanTag 设置tag
func (ot *OpenTelemetryTrace) SetSpanTag(span interface{}, tags map[string]string) {
	for k, v := range tags {
		span.(trace.Span).SetAttributes(attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}
}

// InjectHTTP http传id
func (ot *OpenTelemetryTrace) InjectHTTP(c context.Context, span interface{}, header *http.Header) error {
	header.Set(defaultTraceId, span.(trace.Span).SpanContext().TraceID().String())
	header.Set(defaultSpanId, span.(trace.Span).SpanContext().SpanID().String())
	p := otel.GetTextMapPropagator()
	p.Inject(c, propagation.HeaderCarrier(*header))
	return nil
}

// ExtractHTTP 解析id
func (ot *OpenTelemetryTrace) ExtractHTTP(c context.Context, spanName string, header http.Header) (context.Context, interface{}, error) {
	propagator := otel.GetTextMapPropagator()
	pCtx := propagator.Extract(c, propagation.HeaderCarrier(header))
	traceID := header.Get(traceId)
	spanID := header.Get(spanId)
	//解析出id
	sid, err := trace.SpanIDFromHex(spanID)
	if err != nil {
		return c, nil, err
	}
	tid, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return c, nil, err
	}
	//获取上下文
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
		TraceState: trace.TraceState{},
		Remote:     true,
	})
	sCtx := trace.ContextWithRemoteSpanContext(pCtx, spanCtx)
	ctx, span := otel.Tracer(spanName).Start(sCtx, spanName)
	return ctx, span, nil
}

// ExtractText 解析id key:header的key val:header对应key的值
func (ot *OpenTelemetryTrace) ExtractText(spanName, key, val string) (interface{}, error) {
	return nil, nil
}

// RunHTTP 启动http服务
func (ot *OpenTelemetryTrace) RunHTTP(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

// OnError 设置错误
func (ot *OpenTelemetryTrace) OnError(span interface{}, err error) {
	span.(trace.Span).SetAttributes(attribute.KeyValue{
		Key:   "error",
		Value: attribute.BoolValue(true),
	})
	span.(trace.Span).AddEvent("error", trace.WithAttributes(attribute.KeyValue{
		Key:   "error",
		Value: attribute.StringValue(err.Error()),
	}))
}

// SetSpanIDAndTraceID 设置获取header的参数名
func SetSpanIDAndTraceID(traceIdName, spanIdName string) {
	traceId = traceIdName
	spanId = spanIdName
}
