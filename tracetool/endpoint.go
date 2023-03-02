package tracetool

type traceEndpoint struct {
	f func(*traceConfig)
}

func (te traceEndpoint) apply() {
	te.f(traceConf)
}

func newTraceEndpoint(f func(traceConf *traceConfig)) FunTraceOption {
	return traceEndpoint{f: f}
}

func WithTraceEndpoint(endpoint string) FunTraceOption {
	return newTraceEndpoint(func(traceConf *traceConfig) {
		traceConf.endpoint = endpoint
	})
}
