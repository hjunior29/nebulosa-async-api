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

		apiRouter.POST("/task", task.Create)
		apiRouter.GET("/task/:id", task.Read)
		apiRouter.PUT("/task/:id", task.Update)
		apiRouter.DELETE("/task/:id", task.Delete)

		apiRouter.GET("/ping", status.Ping)
	}
}
