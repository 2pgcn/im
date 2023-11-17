package trace_conf

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"sync"
	"time"
)

var WindowsRate *windowsRate
var WindowsRateOnce sync.Once

const windowsLen = 8

type windowsRate struct {
	mu      sync.Mutex
	windows []int
}

// TraceIDRatioBased 当前请求小于maxRate,返回AlwaysSample,到大最大minRate了后,按fraction比例记录,-1为不限制全按比例
func TraceIDMinNumAndRate(fraction float64, minRepNum int) sdktrace.Sampler {
	if minRepNum == -1 {
		return sdktrace.TraceIDRatioBased(fraction)
	}
	WindowsRateOnce.Do(func() {
		WindowsRate = &windowsRate{
			mu:      sync.Mutex{},
			windows: make([]int, windowsLen),
		}
	})
	winIndex := time.Now().Second() % windowsLen

	WindowsRate.mu.Lock()
	WindowsRate.windows[winIndex]++
	nums := WindowsRate.windows[winIndex]
	defer WindowsRate.mu.Unlock()

	if nums < minRepNum {
		return sdktrace.AlwaysSample()
	}
	return sdktrace.TraceIDRatioBased(fraction)
}

// 继承%采样,保证每秒最低量
//func (as gameimSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
//	result := sdktrace.SamplingResult{
//		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
//	}
//	shouldSample := true
//	winIndex := time.Now().Second() % windowsLen
//	WindowsRate.windows[winIndex]++
//	nums := WindowsRate.windows[winIndex]
//	WindowsRate.mu.Lock()
//	defer WindowsRate.mu.Unlock()
//
//	if nums > as.maxRate {
//		shouldSample = false
//	}
//	if shouldSample {
//		result.Decision = sdktrace.RecordAndSample
//	} else {
//		result.Decision = sdktrace.Drop
//	}
//	return result
//}
//
//func (as gameimSampler) Description() string {
//	return "GameimSampler"
//}
