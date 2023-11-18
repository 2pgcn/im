package server

import (
	"context"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	_ "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	kuberegistry "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlpTraceGrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	slsProjectHeader         = "x-sls-otel-project"
	slsInstanceIDHeader      = "x-sls-otel-instance-id"
	slsAccessKeyIDHeader     = "x-sls-otel-ak-id"
	slsAccessKeySecretHeader = "x-sls-otel-ak-secret"
	slsSecurityTokenHeader   = "x-sls-otel-token"
)

type OtherServer struct {
	tp *tracesdk.TracerProvider
	r  registry.Registrar
}

func (r *OtherServer) Start(ctx context.Context) error {
	tp, err := trace_conf.GetTracerProvider()
	if err != nil {
		return err
	}
	gamelog.Debug("trace_conf start TracerProvider and set global")
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	return nil
}

func (r *OtherServer) Stop(ctx context.Context) error {
	return nil
}

func NewOtherServer() *OtherServer {
	tp, err := getTracerProvider()
	if err != nil {
		panic(err)
	}
	return &OtherServer{
		tp: tp,
	}
}

func getExporterSls() (*otlptrace.Exporter, error) {
	headers := map[string]string{
		slsProjectHeader:         "pg-gameim",
		slsInstanceIDHeader:      "pg-gameim",
		slsAccessKeyIDHeader:     "LTAI5t9xdeJpqKqTapDfUctd",
		slsAccessKeySecretHeader: "qH5IWvC7JQzPG5lIdfhq8WzRcGUHEY",
	}
	traceSecureOption := otlpTraceGrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	traceExporter, err := otlptrace.New(context.Background(),
		otlpTraceGrpc.NewClient(otlpTraceGrpc.WithEndpoint("pg-gameim.cn-shenzhen.log.aliyuncs.com:10010"),
			traceSecureOption,
			otlpTraceGrpc.WithHeaders(headers),
			otlpTraceGrpc.WithCompressor(gzip.Name)))
	return traceExporter, err
}
func getTracerProvider() (*tracesdk.TracerProvider, error) {
	exp, err := getExporterSls()
	//stdouttrace.New()
	//getExporterSls()
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
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

func RegisterK8s(c *conf.Server) registry.Registrar {
	client, err := getK8sClient(c.GetRegister().KubeConfig)
	if err != nil {
		panic(err)
	}
	return kuberegistry.NewRegistry(client)
}

func getK8sClient(c *conf.Kubernetes) (client *kubernetes.Clientset, err error) {
	var restConfig *rest.Config
	if c.KubernetesClientType == conf.KubernetesClientType_INCluster {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", c.KubeConfigPath)

	}
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
