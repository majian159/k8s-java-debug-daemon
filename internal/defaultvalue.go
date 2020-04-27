package internal

import (
	"fmt"
	"javaDebugDaemon/internal/app/nodelock"
	"javaDebugDaemon/internal/app/stackstorage"
	"javaDebugDaemon/internal/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var kubernetesClient *util.KubernetesClient

func DefaultKubernetesClient() (*util.KubernetesClient, error) {
	if kubernetesClient != nil {
		return kubernetesClient, nil
	}

	config, err := getConfigByInCluster()

	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, fmt.Errorf("new k8s client instance error: %v", err)
	}

	kubernetesClient = &util.KubernetesClient{ClientSet: clientSet, Config: config}
	return kubernetesClient, nil
}

func getConfigByInCluster() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getConfigByOutOfCluster() (*rest.Config, error) {
	configFile := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", configFile)

	if err != nil {
		return nil, fmt.Errorf("build config error: %v", err)
	}

	return config, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

var stackStorage stackstorage.StackStorage = stackstorage.NewFileStackStorage()

func DefaultStackStorage() stackstorage.StackStorage {
	return stackStorage
}

var defaultNodeLockManager = nodelock.NewLockManager(10)

func GetDefaultNodeLockManager() *nodelock.LockManager {
	return &defaultNodeLockManager
}
