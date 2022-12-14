package main

import (
	"github.com/manureddy7143/GolangStarter/source/controller"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title Articles Microservice
// @version 0.0.2
// @description This microservice serves as an example microservice
func setupRoutes(r *gin.Engine) {
	// Instantiate controllers
	articleController := controller.UserManagement{}

	// Set application context in URL - Do not edit this
	application := r.Group(viper.GetString("server.basepath"))
	{
		// All routes for API version V1
		v1 := application.Group("/auth")
		{
			docs := v1.Group("/docs")
			{
				docs.StaticFile("/swagger.json", "./docs/swagger.json")
				docs.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
			}

			articles := v1.Group("/users")
			{
				articles.POST(controller.RegisterPath, articleController.Register)
				articles.POST(controller.LoginPath, articleController.Login)
				articles.POST(controller.ProfilePath, articleController.Profile)
			}
		}
	}
}
