package server

import (
	"context"
	"github.com/2pgcn/gameim/conf"
	_ "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	kuberegistry "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	"github.com/go-kratos/kratos/v2/registry"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
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

	return nil
}

func (r *OtherServer) Stop(ctx context.Context) error {
	return nil
}

func NewOtherServer() *OtherServer {
	return &OtherServer{}
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
