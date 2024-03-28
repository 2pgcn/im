package event

import (
	"github.com/2pgcn/gameim/pkg/gamelog"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*queueMsg)(nil)

func (m *queueMsg) Get(key string) string {
	gamelog.Debug("Msg TextMapCarrier: get", key)
	if v, ok := m.H[key]; ok {
		return v.(string)
	}
	return ""
}

func (m *queueMsg) Set(key string, value string) {
	gamelog.Debug("Msg TextMapCarrier: set", key, value)
	m.H[key] = value
}

func (m *queueMsg) Keys() (res []string) {
	e := m.Header()
	for k, _ := range e {
		res = append(res, k)
	}
	return
}
