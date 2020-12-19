package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// User is used by pop to map your users database table to your go code.
type User struct {
	gorm.Model
	ID                string    `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Provider          string    `json:"provider" db:"provider"`
	ProviderID        string    `json:"provider_id" db:"provider_id"`
	TransactionHashes string    `form:"transaction_hashes" json:"transaction_hashes" db:"transaction_hashes"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

func main() {
	dsn := "root:obatala88@tcp(127.0.0.1:3306)/accord365_development?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

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
		fmt.Println("in new payment get")

		// // Create a new user
		// u2 := User{ID: "3", Name: "Andrea Henderson", Provider: "mySite", ProviderID: "4"}
		// result := db.Create(&u2)
		// fmt.Println("Result ", result)
		// Get input from the transaction (txHash)

		c.HTML(http.StatusOK, "new_payment.tmpl", gin.H{
			"title": "New Payment",
		})
	})

	router.POST("/new_payment", func(c *gin.Context) {
		var transaction string
		transaction = c.PostForm("txValueHidden")
		fmt.Printf("tx: %s\n;", transaction)

		var result string
		db.Raw("SELECT transaction_hashes FROM users WHERE id = ?", "2").Scan(&result)
		fmt.Printf("result: %s\n;", result)

		var newTransactionString string
		newTransactionString = (result + "," + transaction)
		fmt.Printf("new tx string: %s\n;", newTransactionString)
		db.Exec("UPDATE users SET transaction_hashes = ? WHERE id = ? ", newTransactionString, "2")

		db.AutoMigrate(&User{})
		db.Save(&User{})

		var displayTransactionHashes string
		db.Raw("SELECT transaction_hashes FROM users WHERE id = ?", "2").Scan(&displayTransactionHashes)
		fmt.Printf("displayTransactionHashes: %s\n;", displayTransactionHashes)
		splitDisplayTransactionHashes := strings.Split(displayTransactionHashes, ",")
		// TODO: Update the splitDisplayTransactions array to remove "[]"
		fmt.Printf("splitDisplayTransactionHashes: %s\n;", splitDisplayTransactionHashes)

		c.HTML(http.StatusOK, "new_payment.tmpl", gin.H{
			"title":                         "New Payment",
			"splitDisplayTransactionHashes": splitDisplayTransactionHashes,
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
