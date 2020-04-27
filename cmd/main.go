package main

import (
	"javaDebugDaemon/internal/app/handler/grafana"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	handler, err := grafana.NewAlertHookHandler()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.POST("/hooks", handler)
	router.StaticFS("/", http.Dir("stacks"))
	err = router.Run()

	if err != nil {
		panic(err)
	}
}
