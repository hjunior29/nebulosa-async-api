package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/api"
	"github.com/hjunior29/nebulosa-async-api/internal/config"
	"github.com/hjunior29/nebulosa-async-api/internal/config/database"
	"github.com/hjunior29/nebulosa-async-api/internal/domain"
	"github.com/hjunior29/nebulosa-async-api/internal/modules/worker"
)

func init() {

	err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate(
		&domain.User{},
		&domain.Task{},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Seed()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	go worker.StartWorker()

	router := gin.Default()
	api.Routes(router)

	log.Fatal(router.Run(":" + config.PORT))
}
