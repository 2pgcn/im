package event

import (
	"github.com/2pgcn/gameim/pkg/gamelog"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*QueueMsg)(nil)

func (m *QueueMsg) Get(key string) string {
	gamelog.GetGlobalog().Debug("Msg TextMapCarrier: get", key)
	if v, ok := m.H[key]; ok {
		return v
	}
	return ""
}

func (m *QueueMsg) Set(key string, value string) {
	gamelog.GetGlobalog().Debug("Msg TextMapCarrier: set", key, value)
	m.H[key] = value
}

func (m *QueueMsg) Keys() (res []string) {
	e := m.Header()
	for k, _ := range *e {
		res = append(res, k)
	}
	return
}
