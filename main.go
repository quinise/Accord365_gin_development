package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Accord365",
		})
	})

	router.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.tmpl", gin.H{
			"title": "Accord365",
		})
	})

	router.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
			"title": "Dashboard",
		})
	})

	router.GET("/wallet_get", func(c *gin.Context) {
		c.HTML(http.StatusOK, "wallet.tmpl", gin.H{
			"title": "Wallet",
		})
	})

	router.POST("/wallet", func(c *gin.Context) {
		c.HTML(http.StatusOK, "wallet.tmpl", gin.H{
			"title": "Wallet",
		})
	})

	router.GET("/new_payment_get", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_payment.tmpl", gin.H{
			"title": "New Payment",
		})
	})

	router.POST("/new_payment", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_payment.tmpl", gin.H{
			"title": "New Payment",
		})
	})

	router.GET("/transaction_get", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transaction.tmpl", gin.H{
			"title": "Transaction History",
		})
	})

	router.POST("/transaction_history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transaction.tmpl", gin.H{
			"title": "Transaction History",
		})
	})

	router.GET("/new_contract_get", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_contract.tmpl", gin.H{
			"title": "New Contract",
		})
	})

	router.POST("/new_contract", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_contract.tmpl", gin.H{
			"title": "New Contract",
		})
	})

	router.GET("/payment_schedule", func(c *gin.Context) {
		c.HTML(http.StatusOK, "payment_schedule.tmpl", gin.H{
			"title": "Payment Schedule",
		})
	})

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")

}
