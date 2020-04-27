package grafana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"javaDebugDaemon/internal"
	"javaDebugDaemon/internal/app"
	"javaDebugDaemon/internal/app/stackstorage"
	"javaDebugDaemon/internal/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var kubernetesClient *util.KubernetesClient
var stackStorage stackstorage.StackStorage = stackstorage.NewFileStackStorage()

func NewAlertHookHandler() (func(ctx *gin.Context), error) {
	client, err := internal.DefaultKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("create k8s client error: %v", err)
	}

	if client == nil {
		log.Printf("get k8s client return nil: %v\n", err)
		return nil, fmt.Errorf("k8s client is nil")
	}

	kubernetesClient = client
	return HookHandler, nil
}

func HookHandler(c *gin.Context) {
	body := c.Request.Body
	if body == nil {
		c.String(http.StatusBadRequest, "not found body.")
		return
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)

	if err != nil {
		log.Printf("read all error: %v\n", err)
		c.String(http.StatusBadRequest, "read body error.")
		return
	}

	if len(data) <= 0 {
		c.String(http.StatusBadRequest, "body is empty.")
		return
	}

	fmt.Println(string(data))

	var model AlertModel
	err = json.Unmarshal(data, &model)

	if err != nil {
		log.Printf("unmarshal error: %v\n", err)
		c.String(http.StatusBadRequest, "unmarshal error")
		return
	}

	go doHandle(model)

	c.Status(http.StatusAccepted)
}

func doHandle(model AlertModel) {
	// 忽略正常告警
	if model.IsOk() {
		return
	}

	matches := model.EvalMatches
	for _, e := range matches {
		go doHandleEvalMatch(e)
	}
}

func doHandleEvalMatch(model EvalMatchModel) {
	tag := model.Tags
	node := tag.Node

	nodeLockManager := *internal.GetDefaultNodeLockManager()
	locker := nodeLockManager.GetLock(node)

	locker.Lock()
	defer locker.Unlock()

	namespace, podName, containerName := tag.Namespace, tag.Pod, tag.Container
	stack, err := app.CrawlString(*kubernetesClient, app.CrawlContext{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Node:          node,
	})

	if err != nil {
		log.Printf("crawString error: %v\n", err)
		return
	}

	err = storeStack(stack, model)

	if err != nil {
		log.Printf("store stack error: stack:%s, %v\n", stack, err)
	}

}

func storeStack(stack string, model EvalMatchModel) error {
	tags := model.Tags
	namespace, podName, containerName, node := tags.Namespace, tags.Pod, tags.Container, tags.Node
	err := stackStorage.Store(stackstorage.ContainerStackModel{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		Node:          node,
		Stack:         stack,
	})

	return err
}
