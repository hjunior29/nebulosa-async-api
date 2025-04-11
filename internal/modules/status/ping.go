package status

import (
	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
)

func Ping(c *gin.Context) {
	utils.SuccessResponse(c, 200, "pong", nil)
}
