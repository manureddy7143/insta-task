package main

import (
	"github.com/manureddy7143/insta-task/source/controller"

	"github.com/gin-gonic/gin"
)

// @title Articles Microservice
// @version 0.0.2
// @description This microservice serves as an example microservice
func setupRoutes(r *gin.Engine) {
	// Instantiate controllers
	transactionController := controller.Transaction{}

	r.POST(controller.Transactions, transactionController.PostTransactions)
	r.GET(controller.GetStat, transactionController.GetStatstics)
	r.DELETE(controller.Transactions, transactionController.DeleteAllTransactions)
	r.POST(controller.Location, transactionController.SetLocation)
	r.DELETE(controller.Location, transactionController.ResetLocation)
}
