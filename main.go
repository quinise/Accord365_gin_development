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
	//"strings"
	"time"

	// "github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"

	// "github.com/accord365/middleware"
	"github.com/gin-gonic/gin"
	// "github.com/gorilla/sessions"
	"github.com/markbates/goth"

	// "github.com/markbates/goth/gothic"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/providers/google"
	"github.com/oov/gothic"
	csrf "github.com/utrack/gin-csrf"

	"golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"

	//"gorm.io/driver/mysql"
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

var cred Credentials
var conf *oauth2.Config

// var state string

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// // SessionUser is a retrieved and authentiacted user.
// type SessionUser struct {
// 	Sub           string `json:"sub"`
// 	Name          string `json:"name"`
// 	GivenName     string `json:"given_name"`
// 	FamilyName    string `json:"family_name"`
// 	Profile       string `json:"profile"`
// 	Picture       string `json:"picture"`
// 	Email         string `json:"email"`
// 	EmailVerified bool   `json:"email_verified"`
// 	Gender        string `json:"gender"`
// }

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
	//sessionSecret := goDotEnvVariable("sessionSecret")

	// key := sessionSecret // Replace with your SESSION_SECRET or similar
	// maxAge := 86400 * 30     // 30 days
	// isProd := false          // Set to true when serving over https

	// store := sessions.NewCookieStore([]byte(key))
	// store.MaxAge(maxAge)
	// store.Options.Path = "/"
	// store.Options.HttpOnly = true // HttpOnly should always be enabled
	// store.Options.Secure = isProd

	clientId := goDotEnvVariable("clientId")
	clientSecret := goDotEnvVariable("clientSecret")

	googleProvider := google.New(clientId, clientSecret, "http://localhost:8080/auth/google/callback")
	goth.UseProviders(googleProvider)

	// dsn := "root:obatala88@tcp(127.0.0.1:3306)/accord365_development?charset=utf8&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()
	// token := randToken()
	// if err != nil {
	// 	log.Fatal("unable to generate random token: ", err)
	// }
	// store := cookie.NewStore([]byte(token))
	// store.Options(sessions.Options{
	// 	Path:   "/",
	// 	MaxAge: 86400 * 7,
	// })
	//router.Use(sessions.Sessions("accord365Session", store))
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "not_found.html", gin.H{
			"title": "Page not found",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Sign In",
			// "user":  userID,
		})

	})

	router.GET("/auth/:google/callback", func(c *gin.Context) {
		user, err := gothic.CompleteAuth(c.Param("google"), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard",
			"user":  user,
		})
	})

	router.GET("/auth/:google", func(c *gin.Context) {
		err := gothic.BeginAuth(c.Param("google"), c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	authorized := router.Group("/auth")
	authorized.GET("/dashboard", func(c *gin.Context) {
		// session := sessions.Default(c)
		// userID := session.Get("userId")
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard",
			// "user":  userID,
		})
	})

	authorized.POST("/wallet", func(c *gin.Context) {

		c.Redirect(http.StatusFound, "/auth/wallet_get")
	})

	authorized.GET("/wallet_get", func(c *gin.Context) {
		// session := sessions.Default(c)
		// userID := session.Get("userId")

		c.HTML(http.StatusOK, "wallet.html", gin.H{
			"title": "Wallet",
			// "user":  userID,
		})
	})

	// authorized.POST("/new-payment", func(c *gin.Context) {
	// 	fmt.Println("*****In new-payment")
	// 	// TODO (Production): update database to permit unlimited char/VAR for transaction_hashes
	// 	var transaction string
	// 	transaction = c.PostForm("txValueHidden")
	// 	if len(transaction) < 1 {
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else {
	// 		fmt.Printf("tx: %s\n;", transaction)
	// 	}

	// 	var result string
	// 	if len(transaction) < 1 {
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else if len(transaction) > 1 {
	// 		db.Raw("SELECT transaction_hashes FROM users WHERE id = ?", "2").Scan(&result)
	// 		fmt.Printf("result: %s\n;", result)
	// 	}

	// 	var newTransactionString string
	// 	if result == "null" || len(result) < 1 {
	// 		newTransactionString = transaction
	// 		fmt.Printf("new transaction string (no new entries): %s\n;", newTransactionString)
	// 		db.Exec("UPDATE users SET transaction_hashes = ? WHERE id = ? ", newTransactionString, "2")
	// 	} else if len(result) > 1 {
	// 		newTransactionString = (result + "," + transaction)
	// 		fmt.Printf("new tx string: %s\n;", newTransactionString)
	// 		db.Exec("UPDATE users SET transaction_hashes = ? WHERE id = ? ", newTransactionString, "2")
	// 	}

	// 	var displayTransactionHashes string
	// 	var splitDisplayTransactionHashes []string
	// 	if len(newTransactionString) < 1 {
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else if len(newTransactionString) > 1 {
	// 		db.Raw("SELECT transaction_hashes FROM users WHERE id = ?", "2").Scan(&displayTransactionHashes)
	// 		fmt.Printf("displayTransactionHashes: %s\n;", displayTransactionHashes)
	// 		splitDisplayTransactionHashes = strings.Split(displayTransactionHashes, ",")
	// 		fmt.Printf("splitDisplayTransactionHashes: %s\n;", splitDisplayTransactionHashes)
	// 	}

	// 	session := sessions.Default(c)
	// 	var splitDisplayTransactionHashesStr string
	// 	if len(displayTransactionHashes) < 1 {
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else if len(displayTransactionHashes) > 1 {
	// 		splitDisplayTransactionHashesStr = strings.Join(splitDisplayTransactionHashes, "\n")
	// 		session.Set("splitDisplayTransactionHashesString", splitDisplayTransactionHashesStr)
	// 		session.Save()
	// 		fmt.Printf("splitDisplayTransactionHashesStr %s\n", splitDisplayTransactionHashesStr)
	// 	}

	// 	if len(splitDisplayTransactionHashesStr) > 1 {
	// 		db.AutoMigrate(&User{})
	// 		db.Save(&User{})
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else if len(splitDisplayTransactionHashesStr) < 1 {
	// 		c.Redirect(http.StatusFound, "/index")
	// 	}
	// })

	// authorized.GET("/new_payment_get", func(c *gin.Context) {
	// 	fmt.Println("%%%%c.Request GET", c.Request)

	// 	session := sessions.Default(c)
	// 	userID := session.Get("userId")
	// 	session.Delete("_csrf")
	// 	session.Set("_csrf", csrf.GetToken(c))
	// 	session.Save()
	// 	_csrf := session.Get("_csrf")

	// 	splitDisplayTransactionHashesString := session.Get("splitDisplayTransactionHashesString")
	// 	if splitDisplayTransactionHashesString != session.Get("splitDisplayTransactionHashesString") {
	// 		fmt.Println("Error: Issue with retrieving transaction string.")
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else {
	// 		fmt.Printf("Session test: %s\n;", splitDisplayTransactionHashesString)
	// 	}

	// 	// Set http request header if csrf token is ok
	// 	if _csrf, ok := session.Get("_csrf").(string); !ok {
	// 		fmt.Println("@@@@@Find Request header FAILED", c.Request.Header.Values("X-CSRF-Token"))
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else {
	// 		c.Request.Header.Set("X-CSRF-Token", _csrf)
	// 		fmt.Println("@@@@@Request header SUCCESS", c.Request.Header.Values("X-CSRF-Token"))
	// 	}

	// 	header := c.Request.Header.Values("X-CSRF-Token")
	// 	// compare csfr token to header
	// 	if _csrf != header[0] {
	// 		fmt.Println("***csrf tokens don't match!*** " + _csrf.(string) + " " + header[0])
	// 		c.Redirect(http.StatusFound, "/auth/new_payment_get")
	// 	} else if _csrf == header[0] {
	// 		fmt.Println("***csrf tokens match!***")

	// 		c.HTML(http.StatusOK, "new_payment.html", gin.H{
	// 			"_csrf": _csrf,
	// 			"title": "New Payment",
	// 			"user":  userID,
	// 			// "splitDisplayTransactionHashesString": splitDisplayTransactionHashesString,
	// 		})

	// 	}

	// })

	authorized.POST("/transaction_history", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth/transaction_get")
	})

	// authorized.GET("/transaction_get", func(c *gin.Context) {
	// 	session := sessions.Default(c)
	// 	userID := session.Get("userId")
	// 	session.Set("csrfFieldTransaction", csrf.GetToken(c))
	// 	csrfFieldTransaction := session.Get("csrfFieldTransaction")

	// 	c.HTML(http.StatusOK, "transaction.html", gin.H{
	// 		"csrfFieldTransaction": csrfFieldTransaction,
	// 		"title":                "Transaction History",
	// 		"user":                 userID,
	// 	})
	// })

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
		// session := sessions.Default(c)
		// userID := session.Get("userId")

		c.HTML(http.StatusOK, "new_contract.html", gin.H{
			"csrfFieldChild":     csrf.GetToken(c),
			"csrfFieldProbation": csrf.GetToken(c),
			"csrfFieldTraffic":   csrf.GetToken(c),
			"title":              "New Contract",
			// "user":               userID,
		})
	})

	authorized.GET("/logout", func(c *gin.Context) {
		// session := sessions.Default(c)
		// session.Clear()
		// session.Save()
		c.Redirect(http.StatusFound, "/index")
	})

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
