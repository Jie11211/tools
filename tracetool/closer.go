package tracetool

type openTelemetryClose struct {
	f func()
}

func (otlpc *openTelemetryClose) Close() error {
	otlpc.f()
	return nil
}
