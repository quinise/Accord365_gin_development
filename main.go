package main

import (
	// "crypto/rand"
	// "encoding/base64"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "log"
	"net/http"
	// "os"
	"strings"
	"time"

	// "github.com/accord365/middleware"
	// "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var cred Credentials
// var conf *oauth2.Config

var db *gorm.DB
var err error

// // Credentials which stores google ids.
// type Credentials struct {
// 	Cid     string `json:"cid"`
// 	Csecret string `json:"csecret"`
// }

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

// func init() {
// 	file, err := ioutil.ReadFile("./creds.json")
// 	if err != nil {
// 		log.Printf("File error: %v\n", err)
// 		os.Exit(1)
// 	}
// 	if err := json.Unmarshal(file, &cred); err != nil {
// 		log.Println("unable to marshal data")
// 		return
// 	}

// 	conf = &oauth2.Config{
// 		ClientID:     cred.Cid,
// 		ClientSecret: cred.Csecret,
// 		RedirectURL:  "http://127.0.0.1:8080/authorize",
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
// 		},
// 		Endpoint: google.Endpoint,
// 	}
// }

// // RandToken generates a random @l length token.
// func RandToken(l int) (string, error) {
// 	b := make([]byte, l)
// 	if _, err := rand.Read(b); err != nil {
// 		return "", err
// 	}
// 	return base64.StdEncoding.EncodeToString(b), nil
// }

// func getLoginURL(state string) string {
// 	return conf.AuthCodeURL(state)
// }

func main() {
	dsn := "root:obatala88@tcp(127.0.0.1:3306)/accord365_development?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()
	// token, err := RandToken(64)
	// if err != nil {
	// 	log.Fatal("unable to generate random token: ", err)
	// }
	// store := sessions.NewCookieStore([]byte(token))
	// store.Options(sessions.Options{
	// 	Path:   "/",
	// 	MaxAge: 86400 * 7,
	// })

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(sessions.Sessions("mysession", store))
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	// // authorization handles authentication of a user and initiates a session.
	// router.GET("/authorization", func(c *gin.Context) {
	// 	// Handle the exchange code to initiate a transport.
	// 	session := sessions.Default(c)
	// 	retrievedState := session.Get("state")
	// 	queryState := c.Request.URL.Query().Get("state")
	// 	if retrievedState != queryState {
	// 		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
	// 		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
	// 		return
	// 	}
	// 	code := c.Request.URL.Query().Get("code")
	// 	tok, err := conf.Exchange(oauth2.NoContext, code)
	// 	if err != nil {
	// 		log.Println(err)
	// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
	// 		return
	// 	}

	// 	client := conf.Client(oauth2.NoContext, tok)
	// 	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	// 	if err != nil {
	// 		log.Println(err)
	// 		c.AbortWithStatus(http.StatusBadRequest)
	// 		return
	// 	}
	// 	defer userinfo.Body.Close()
	// 	data, _ := ioutil.ReadAll(userinfo.Body)
	// 	u := &User{}
	// 	if err = json.Unmarshal(data, &u); err != nil {
	// 		log.Println(err)
	// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
	// 		return
	// 	}
	// 	session.Set("user-id", u.ID)
	// 	err = session.Save()
	// 	if err != nil {
	// 		log.Println(err)
	// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
	// 		return
	// 	}

	// 	result := db.First(&u.ID)
	// 	if result.Error == nil {
	// 		log.Println(err)
	// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
	// 	} else {
	// 		db.Save(&u)
	// 	}

	// 	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{})
	// })

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Accord365",
		})
	})

	router.GET("/signin", func(c *gin.Context) {
		// state, err := RandToken(32)
		// fmt.Printf("state: %s", state)
		// if err != nil {
		// 	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while generating random data."})
		// 	return
		// }
		// session := sessions.Default(c)
		// log.Println("session: ", session)
		// session.Set("state", state)
		// err = session.Save()
		// if err != nil {
		// 	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		// 	return
		// }

		// link := getLoginURL(state)

		c.HTML(http.StatusOK, "signin.tmpl", gin.H{
			"title": "Accord365",
			// "link":  link,
		})
	})

	// authorize.Use(middleware.AuthorizeRequest())
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

		splitDisplayTransactionHashesStr := strings.Join(splitDisplayTransactionHashes, "\n")
		c.Set("splitDisplayTransactionHashesStr", splitDisplayTransactionHashesStr)
		fmt.Printf("splitDisplayTransactionHashesStr %s\n", splitDisplayTransactionHashesStr)

		c.HTML(http.StatusOK, "new_payment.tmpl", gin.H{
			"title":                            "New Payment",
			"splitDisplayTransactionHashesStr": splitDisplayTransactionHashesStr,
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
		c.HTML(http.StatusOK, "new_contract.html", gin.H{
			"title": "New Contract",
		})
	})

	router.POST("/new_contract", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new_contract.html", gin.H{
			"title": "New Contract",
		})
	})

	router.GET("/payment_schedule", func(c *gin.Context) {
		c.HTML(http.StatusOK, "payment_schedule.tmpl", gin.H{
			"title": "Payment Schedule",
		})
	})

	//create a logout function

	// Listen and server on 0.0.0.0:8081
	router.Run(":8081")
}
