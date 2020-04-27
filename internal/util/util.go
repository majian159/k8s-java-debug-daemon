package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/kubernetes"
)

type KubernetesClient struct {
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
}

func DefaultKubernetesClient() (*KubernetesClient, error) {
	kubeConfig := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)

	if err != nil {
		return nil, fmt.Errorf("build config error: %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, fmt.Errorf("new k8s client instance error: %v", err)
	}

	return &KubernetesClient{ClientSet: clientSet, Config: config}, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func (client KubernetesClient) Exec(namespace, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer) ([]byte, error) {
	clientSet, config := client.ClientSet, client.Config

	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(
			&v12.PodExecOptions{
				Container: containerName,
				Command:   command,
				Stdin:     stdin != nil,
				Stdout:    stdout != nil,
				Stderr:    true,
				TTY:       false,
			}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	if err != nil {
		return nil, fmt.Errorf("error while creating Executor: %v", err)
	}

	var stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		return stderr.Bytes(), fmt.Errorf("error in Stream: %v", err)
	}

	return stderr.Bytes(), nil
}
