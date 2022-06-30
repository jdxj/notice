package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jdxj/notice/logger"
)

var (
	server *http.Server
)

func Start() {
	r := gin.Default()
	register(r)

	server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Errorf("listen err: %s", err)
			}
		}
	}()
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("shutdown err: %s", err)
	}
}

func register(root gin.IRouter) {
	root.GET("/", hello)
}
