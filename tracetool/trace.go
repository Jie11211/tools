package tracetool

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	Jaeger           = "jaeger"
	OpenTelemetry    = "OpenTelemetry"
	jaegerUrl        = "/api/traces"
	openTelemetryUrl = "/api/otlp/traces"
)

var (
	traceConf = &traceConfig{}
	closer    io.Closer
	url       = ""
	urlMap    = map[string]string{
		Jaeger:        jaegerUrl,
		OpenTelemetry: openTelemetryUrl,
	}
)

type traceConfig struct {
	endpoint string
	url      string
}

type FunTraceOption interface {
	apply()
}

// New tp:使用的链路追踪方式 service:应用名称
func New(tp, service string, opts ...FunTraceOption) (Trace, io.Closer) {
	url = jaegerUrl
	if v, ok := urlMap[tp]; ok {
		url = v
	}
	for _, opt := range opts {
		opt.apply()
	}
	switch tp {
	case OpenTelemetry:
		closer = traceConf.newOpenTelemetryTracer(service)
		return &OpenTelemetryTrace{}, closer
	default:
		jaegerTrace := &JaegerTrace{}
		jaegerTrace.tracer, closer = traceConf.newJaegerTracer(service)
		return jaegerTrace, closer
	}
}

// newJaegerTracer 新建jaeger链路追踪
func (traceConf *traceConfig) newJaegerTracer(service string) (opentracing.Tracer, io.Closer) {
	if !strings.Contains(traceConf.endpoint, "http://") {
		traceConf.endpoint = "http://" + traceConf.endpoint
	}
	sender := transport.NewHTTPTransport(
		fmt.Sprintf("%s%s", traceConf.endpoint, traceConf.url),
	)
	tracer, closer := jaeger.NewTracer(service,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender))
	return tracer, closer
}

// newOpenTelemetryTracer 新建OpenTelemetry链路追踪
func (traceConf *traceConfig) newOpenTelemetryTracer(service string) io.Closer {
	traceConf.endpoint = strings.Replace(traceConf.endpoint, "http://", "", 1)
	traceClientHttp := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(traceConf.endpoint),
		otlptracehttp.WithURLPath(traceConf.url),
		otlptracehttp.WithInsecure())
	otlptracehttp.WithCompression(1)
	exporter, err := otlptrace.New(context.Background(), traceClientHttp)
	if err != nil {
		fmt.Println(err)
	}
	res, err := resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(service),
		),
	)
	if err != nil {
		fmt.Println(err)
	}
	bsp := trace.NewBatchSpanProcessor(exporter)
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(provider)
	openTelemetry := &openTelemetryClose{f: func() {
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := provider.Shutdown(c); err != nil {
			otel.Handle(err)
		}
	}}
	return openTelemetry
}
