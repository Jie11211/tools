package tracetool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestTrace(t *testing.T) {
	tr, closer := New(Jaeger, "test",
		WithTraceEndpoint("http://tracing-analysis-dc-hz.aliyuncs.com"),
		WithTraceUrl(""),
	)
	defer closer.Close()
	ctx, span := tr.Start(context.Background(), "todo1-start")
	defer tr.Finish(span)
	tr.SetSpanTag(span, map[string]string{
		"日期": time.Now().Format("2006-01-02 15:04:05"),
		"id": "55555555555555",
	})
	F1(ctx, tr, span)

	request, err := http.NewRequest(http.MethodGet, "http://localhost:8888/t", nil)
	if err != nil {
		fmt.Println(err)
	}
	tr.InjectHTTP(context.Background(), span, &request.Header)
	fmt.Println(request.Header)
	client := http.Client{}
	do, err := client.Do(request)
	if err != nil {
		tr.OnError(span, err)
	}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		tr.OnError(span, err)
		return
	}
	fmt.Println(string(all))
}

func F1(ctx context.Context, a Trace, span interface{}) {
	_, ctxSpan := a.Transmit(ctx, span, "f1")
	defer a.Finish(ctxSpan)
	a.SetSpanTag(ctxSpan, map[string]string{
		"日期":   time.Now().Format("2006-01-02 15:04:05"),
		"data": "child",
		"id":   "6666666666666666",
	})

}
