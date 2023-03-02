package tracetool

import "fmt"

type traceUrl struct {
	f func(*traceConfig)
}

func (tu traceUrl) apply() {
	tu.f(traceConf)
}

func newTraceUrl(f func(traceConf *traceConfig)) FunTraceOption {
	return traceUrl{f: f}
}

func WithTraceUrl(token string) FunTraceOption {
	return newTraceUrl(func(traceConf *traceConfig) {
		traceConf.url = fmt.Sprintf("/adapt_%s%s", token, url)
	})
}
