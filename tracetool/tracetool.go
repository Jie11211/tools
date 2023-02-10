package tracetool

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	opentrace "go.opentelemetry.io/otel/trace"
)

type TraceProvider struct {
	Tp   *trace.TracerProvider
	Span map[string]*OpenSpan
}

type OpenSpan struct {
	Ctx  context.Context
	Span opentrace.Span
}

func newExporter(url string) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}

// semconv.ServiceNameKey.String(name),
// semconv.ServiceVersionKey.String(version),
// attribute.String("environment", "demo"),
func newResource(attribute ...attribute.KeyValue) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute...,
		),
	)
	return r
}

func NewKeyValue(name, version string, attributeMap map[string]string) (attrKeyValue []attribute.KeyValue) {
	attrKeyValue = append(attrKeyValue, semconv.ServiceNameKey.String(name))
	attrKeyValue = append(attrKeyValue, semconv.ServiceVersionKey.String(version))
	for k, v := range attributeMap {
		attrKeyValue = append(attrKeyValue, attribute.String(k, v))
	}
	return
}

func NewTraceProvider(url string, attribute ...attribute.KeyValue) *TraceProvider {
	exp, _ := newExporter(url)
	// 创建链路生成器,这里将导出器与资源信息配置进去.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource(attribute...)),
	)
	otel.SetTracerProvider(tp)
	return &TraceProvider{Tp: tp}
}

func (tr *TraceProvider) Shutdown(ctx context.Context) error {
	return tr.Tp.Shutdown(ctx)
}

func (tr *TraceProvider) NewSpan(ctx context.Context, tranceName, spanName string) (context.Context, opentrace.Span) {
	startCtx, span := otel.Tracer(tranceName).Start(ctx, spanName)
	if _, ok := tr.Span[spanName]; !ok {
		tr.Span[spanName] = &OpenSpan{
			Span: span,
			Ctx:  startCtx,
		}
		return startCtx, span
	}
	return tr.Span[spanName].Ctx, tr.Span[spanName].Span
}

func (tr *TraceProvider) GetSpanByHttpHeader(c context.Context, header *http.Header, tranceName, spanName, funcName string) (context.Context, opentrace.Span, error) {
	propagator := otel.GetTextMapPropagator()
	pctx := propagator.Extract(c, propagation.HeaderCarrier(*header))
	traceID := header.Get(tranceName)
	spanID := header.Get(spanName)

	//解析出id
	spanid, err := opentrace.SpanIDFromHex(spanID)
	if err != nil {
		return nil, nil, err
	}
	traceid, err := opentrace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, nil, err
	}
	//获取上下文
	spanCtx := opentrace.NewSpanContext(opentrace.SpanContextConfig{
		TraceID:    traceid,
		SpanID:     spanid,
		TraceFlags: opentrace.FlagsSampled, //这个没写，是不会记录的
		TraceState: opentrace.TraceState{},
		Remote:     true,
	})
	sctx := opentrace.ContextWithRemoteSpanContext(pctx, spanCtx)
	ctx, s := tr.Tp.Tracer(funcName).Start(sctx, funcName)
	return ctx, s, nil
}

func (ospan *OpenSpan) AddHeader(request *http.Request, tranceName, spanName string) {
	request.Header.Set(tranceName, ospan.Span.SpanContext().TraceID().String())
	request.Header.Set(spanName, ospan.Span.SpanContext().SpanID().String())
	p := otel.GetTextMapPropagator()
	p.Inject(ospan.Ctx, propagation.HeaderCarrier(request.Header))
}
