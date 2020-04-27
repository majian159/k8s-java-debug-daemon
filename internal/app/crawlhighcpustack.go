package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"javaDebugDaemon/internal/util"
	"strings"
)

type CrawlContext struct {
	Namespace     string
	PodName       string
	ContainerName string
	Node          string
}

const shellFile = "craw.sh"
const targetFile = "/root/craw.sh"

func CrawlString(client util.KubernetesClient, context CrawlContext) (string, error) {
	data, err := crawl(client, context)

	if err != nil {
		return "", err
	}

	str := string(data)

	const startKey = "| plaintext"
	startIndex := strings.Index(str, startKey) + len(startKey) + 2

	str = str[startIndex:]

	const endKey = "[arthas@"
	endIndex := strings.LastIndex(str, endKey) - 2
	str = str[:endIndex]

	return str, nil
}

func crawl(client util.KubernetesClient, context CrawlContext) (stdoutBytes []byte, err error) {
	namespace, podName, containerName, node := context.Namespace, context.PodName, context.ContainerName, context.Node
	_ = node

	commands := []string{"/bin/bash", "-c"}

	commands = append(commands, fmt.Sprintf("cp -f /dev/stdin %[1]s;chmod +x %[1]s;%[1]s", targetFile))

	scriptData, err := ioutil.ReadFile(shellFile)
	if err != nil {
		return nil, fmt.Errorf("read file error: %v", err)
	}

	stdin := bytes.NewReader(scriptData)
	var stdout bytes.Buffer
	stderr, err := client.Exec(namespace, podName, containerName, commands, stdin, &stdout)

	fmt.Println(podName)
	if len(stderr) != 0 {
		return nil, fmt.Errorf("STDERR: " + (string)(stderr))
	}

	if err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}
