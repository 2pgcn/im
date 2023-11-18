package trace_conf

import (
	"context"
	"github.com/2pgcn/gameim/conf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlpTraceGrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/encoding/gzip"
	"sync"
)

const (
	slsProjectHeader         = "x-sls-otel-project"
	slsInstanceIDHeader      = "x-sls-otel-instance-id"
	slsAccessKeyIDHeader     = "x-sls-otel-ak-id"
	slsAccessKeySecretHeader = "x-sls-otel-ak-secret"
	slsSecurityTokenHeader   = "x-sls-otel-token"
)

var traceConfig *conf.TraceConf
var one sync.Once

func SetTraceConfig(traceConf *conf.TraceConf) {
	one.Do(func() {
		traceConfig = traceConf
	})
}

func getExporterSls() (*otlptrace.Exporter, error) {
	headers := map[string]string{
		slsProjectHeader:         traceConfig.SlsProjectHeader,
		slsInstanceIDHeader:      traceConfig.SlsInstanceIDHeader,
		slsAccessKeyIDHeader:     traceConfig.SlsAccessKeyIDHeader,
		slsAccessKeySecretHeader: traceConfig.SlsAccessKeySecretHeader,
	}
	//traceSecureOption := otlpTraceGrpc.WithTLSCredentials(
	//	credentials.NewClientTLSFromCert(nil, ""))

	traceExporter, err := otlptrace.New(context.Background(),
		otlpTraceGrpc.NewClient(
			otlpTraceGrpc.WithEndpoint(traceConfig.Endpoint),
			//traceSecureOption,
			otlpTraceGrpc.WithHeaders(headers),
			otlpTraceGrpc.WithCompressor(gzip.Name)))
	return traceExporter, err
}
func GetTracerProvider() (*tracesdk.TracerProvider, error) {
	exp, err := getExporterSls()
	//"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	//exp, err := stdouttrace.New()
	//stdouttrace.New()
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(TraceIDMinNumAndRate(0.1, 2))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("gameim"),
			attribute.String("env", "dev"),
		)),
	)
	return tp, nil
}

func SetTrace(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("gameim").Start(ctx, spanName, opts...)
	defer span.End()
	return ctx, span
}
