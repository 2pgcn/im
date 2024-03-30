package client

import (
	_ "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	kuberegistry "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/env"
)

var defaultKConfigEnvKey = "KUBECONFIG"

func NewK8sDiscovery(configPath string) (registry.Discovery, error) {
	if len(configPath) == 0 {
		configPath = env.GetString(defaultKConfigEnvKey, "")

	}
	restConfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	kClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return kuberegistry.NewRegistry(kClient), nil
}
