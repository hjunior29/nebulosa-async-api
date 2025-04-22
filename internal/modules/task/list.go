package task

import (
	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/config/database"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
)

func List(c *gin.Context) {
	var tasks []domain.Task
	repository := database.NewRepository(&tasks, c)

	err := repository.FindAllWhere(nil)
	if err != nil {
		utils.ErrorResponse(c, 404, "Error reading tasks", err)
		return
	}

	utils.SuccessResponse(c, 200, "Tasks retrieved successfully", tasks)
}
