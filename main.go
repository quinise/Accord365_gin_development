package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	//"io/ioutil"
	"log"
	//"html/template"
	"net/http"

	"os"
	"path/filepath"

	"strings"
	"time"

	"github.com/accord365/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"

	"github.com/joho/godotenv"
	"github.com/markbates/goth/providers/google"
	"github.com/oov/gothic"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// User is used to map your users database table to your go code.
type User struct {
	gorm.Model
	ID                int       `json:"id" db:"id" gorm:"primaryKey;autoIncrement:true"`
	Name              string    `json:"name" db:"name"`
	Email             string    `json:"email" db:"email"`
	Provider          string    `json:"provider" db:"provider"`
	ProviderID        string    `json:"provider_id" db:"provider_id"`
	TransactionHashes string    `form:"transaction_hashes" json:"transaction_hashes" db:"transaction_hashes"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// RandToken generates a random @l length token.
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	clientId := goDotEnvVariable("clientId")
	clientSecret := goDotEnvVariable("clientSecret")

	googleProvider := google.New(clientId, clientSecret, "http://localhost:8080/auth/google/callback")
	goth.UseProviders(googleProvider)

	dsn := "root:obatala88@tcp(127.0.0.1:3306)/accord365_development?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()

	token := randToken()
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}

	store := cookie.NewStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	router.Use(sessions.Sessions("accord365Session", store))
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "not_found.html", gin.H{
			"title": "Page not found",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Accord365",
		})

	})

	router.GET("/auth/:google", func(c *gin.Context) {
		err := gothic.BeginAuth(c.Param("google"), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	router.GET("/auth/:google/callback", func(c *gin.Context) {
		user, err := gothic.CompleteAuth(c.Param("google"), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			session := sessions.Default(c)
			session.Set("userId", user)
			session.Save()
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard",
			"user":  user,
		})
	})

	authorized := router.Group("/auth")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/dashboard", func(c *gin.Context) {
			c.HTML(http.StatusOK, "dashboard.html", gin.H{
				"title": "Dashboard",
			})
		})
	
		authorized.POST("/wallet", func(c *gin.Context) {
	
			c.Redirect(http.StatusFound, "/auth/wallet_get")
		})
	
		authorized.GET("/wallet_get", func(c *gin.Context) {
			c.HTML(http.StatusOK, "wallet.html", gin.H{
				"title": "Wallet",
			})
		})
	
		authorized.GET("/new_payment_get", func(c *gin.Context) {
			fmt.Println("%%%%c.Request GET", c.Request)
	
			session := sessions.Default(c)
			session.Save()
	
			splitDisplayTransactionHashesString := session.Get("splitDisplayTransactionHashesString")
			if splitDisplayTransactionHashesString != session.Get("splitDisplayTransactionHashesString") {
				fmt.Println("Error: Issue with retrieving transaction string.")
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else {
				fmt.Printf("Session test: %s\n;", splitDisplayTransactionHashesString)
			}
	
			c.HTML(http.StatusOK, "new_payment.html", gin.H{
				"title": "New Payment",
				// "splitDisplayTransactionHashesString": splitDisplayTransactionHashesString,
			})
		})
	
		authorized.POST("/new-payment", func(c *gin.Context) {
			fmt.Println("*****In new-payment")
			// TODO (Production): update database to permit unlimited char/VAR for transaction_hashes
			var transaction string
			transaction = c.PostForm("txValueHidden")
			if len(transaction) < 1 {
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else {
				fmt.Printf("tx: %s\n;", transaction)
			}
	
			var result string
			if len(transaction) < 1 {
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else if len(transaction) > 1 {
				db.Raw("SELECT transaction_hashes FROM users WHERE id = ?", "2").Scan(&result)
				fmt.Printf("result: %s\n;", result)
			}
	
			var newTransactionString string
			if result == "null" || len(result) < 1 {
				newTransactionString = transaction
				fmt.Printf("new transaction string (no new entries): %s\n;", newTransactionString)
				db.Exec("UPDATE users SET transaction_hashes = ? WHERE id = ? ", newTransactionString, "2")
			} else if len(result) > 1 {
				newTransactionString = (result + "," + transaction)
				fmt.Printf("new tx string: %s\n;", newTransactionString)
				db.Exec("UPDATE users SET transaction_hashes = ? WHERE id = ? ", newTransactionString, "2")
			}
	
			var displayTransactionHashes string
			var splitDisplayTransactionHashes []string
			if len(newTransactionString) < 1 {
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else if len(newTransactionString) > 1 {
				db.Raw("SELECT transaction_hashes FROM users WHERE id = ?", "2").Scan(&displayTransactionHashes)
				fmt.Printf("displayTransactionHashes: %s\n;", displayTransactionHashes)
				splitDisplayTransactionHashes = strings.Split(displayTransactionHashes, ",")
				fmt.Printf("splitDisplayTransactionHashes: %s\n;", splitDisplayTransactionHashes)
			}
	
			session := sessions.Default(c)
			var splitDisplayTransactionHashesStr string
			if len(displayTransactionHashes) < 1 {
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else if len(displayTransactionHashes) > 1 {
				splitDisplayTransactionHashesStr = strings.Join(splitDisplayTransactionHashes, "\n")
				session.Set("splitDisplayTransactionHashesString", splitDisplayTransactionHashesStr)
				session.Save()
				fmt.Printf("splitDisplayTransactionHashesStr %s\n", splitDisplayTransactionHashesStr)
			}
	
			if len(splitDisplayTransactionHashesStr) > 1 {
				db.AutoMigrate(&User{})
				db.Save(&User{})
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else if len(splitDisplayTransactionHashesStr) < 1 {
				c.Redirect(http.StatusFound, "/index")
			}
		})
	
		authorized.POST("/transaction_history", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/auth/transaction_get")
		})
	
		authorized.GET("/transaction_get", func(c *gin.Context) {
	
			c.HTML(http.StatusOK, "transaction.html", gin.H{
				"title": "Transaction History",
			})
		})
	
		authorized.POST("/new_contract", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/auth/new_contract_get")
		})
	
		authorized.POST("/upload", func(c *gin.Context) {
			name := c.PostForm("name")
	
			// Source
			file, err := c.FormFile("file")
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
				return
			}
	
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
	
			c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with field name=%s.", file.Filename, name))
			// TODO: update the database to include file in a BLOB
	
			c.Redirect(http.StatusFound, "/auth/new_contract_get")
		})
	
		authorized.GET("/new_contract_get", func(c *gin.Context) {
			c.HTML(http.StatusOK, "new_contract.html", gin.H{
				"title": "New Contract",
			})
		})
	
		authorized.GET("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Clear()
			session.Save()
			c.Redirect(http.StatusFound, "/")
		})
	}

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
