package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/modules/auth"
	"github.com/hjunior29/nebulosa-async-api/internal/modules/status"
	"github.com/hjunior29/nebulosa-async-api/internal/modules/task"
)

func Routes(router *gin.Engine) {
	apiRouter := router.Group("/api")
	{
		apiRouter.POST("/auth/login", auth.Login)

		apiRouter.POST("/task", AuthMiddleware(), task.Create)
		apiRouter.GET("/task/:id", AuthMiddleware(), task.Read)
		apiRouter.GET("/task/:id", AuthMiddleware(), task.Update)
		apiRouter.PUT("/task/:id", AuthMiddleware(), task.Delete)

		apiRouter.GET("/ping", AuthMiddleware(), status.Ping)
	}
}
