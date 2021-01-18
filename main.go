package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/accord365/middleware"
	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// User is used to map your users database table to your go code.
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

var cred Credentials
var conf *oauth2.Config

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// SessionUser is a retrieved and authentiacted user.
type SessionUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &cred); err != nil {
		log.Println("unable to marshal data")
		return
	}

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://127.0.0.1:8081/auth_google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func main() {
	dsn := "root:obatala88@tcp(127.0.0.1:3306)/accord365_development?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()

	token, err := RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("goquestsession", store))
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	router.GET("/index", func(c *gin.Context) {
		fmt.Println("in index")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Accord365",
		})
	})

	router.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", gin.H{
			"title": "Accord365",
		})
	})

	// Google LoginHandler handles the login procedure
	router.GET("/signin_google", func(c *gin.Context) {
		fmt.Println("in signin_google")

		state, err := RandToken(32)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while generating random data."})
			return
		}
		session := sessions.Default(c)
		session.Set("state", state)
		err = session.Save()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
			return
		}
		link := getLoginURL(state)
		c.HTML(http.StatusOK, "auth.html", gin.H{
			"link":  link,
			"title": "Accord365",
		})
	})

	// Google AuthHandler handles authentication of a user and initiates a session.
	router.GET("/auth_google", func(c *gin.Context) {
		fmt.Println("in auth_google")

		// Handle the exchange code to initiate a transport.
		session := sessions.Default(c)
		retrievedState := session.Get("state")
		queryState := c.Request.URL.Query().Get("state")
		if retrievedState != queryState {
			log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
			return
		}
		code := c.Request.URL.Query().Get("code")
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
			return
		}

		client := conf.Client(oauth2.NoContext, tok)
		userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		defer userinfo.Body.Close()
		data, _ := ioutil.ReadAll(userinfo.Body)
		u := SessionUser{}
		if err = json.Unmarshal(data, &u); err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
			return
		}
		session.Set("userId", u.Email)
		err = session.Save()
		log.Println("SessionUsers", u)
		log.Println("Session", session.Get("userId"))

		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
			return
		}
		seen := false
		if session.Get("userId") == u.Email {
			seen = true
		} else {
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving user. Please try again."})
			return
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":  "Dashboard",
			"userId": session.Get("userId"),
			"seen":   seen,
		})
	})

	authorized := router.Group("/auth")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/dashboard", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")
			c.HTML(http.StatusOK, "dashboard.html", gin.H{
				"title": "Dashboard",
				"user":  userID,
			})
		})

		authorized.GET("/wallet_get", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "wallet.html", gin.H{
				"title": "Wallet",
				"user":  userID,
			})
		})

		authorized.POST("/wallet", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "wallet.html", gin.H{
				"title": "Wallet",
				"user":  userID,
			})

		})

		authorized.GET("/new_payment_get", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")
			// splitDisplayTransactionHashesString := c.Get("splitDisplayTransactionHashesStr")

			fmt.Println("in new payment get")

			// TODO: create a new user upon first login
			// Create a new user
			// u2 := User{ID: "2", Name: "Jammal Hendrix", Provider: "mySite", ProviderID: "4"}
			// result := db.Create(&u2)
			// fmt.Println("Result ", result)
			// Get input from the transaction (txHash)
			c.HTML(http.StatusOK, "new_payment.html", gin.H{
				"title": "New Payment",
				"user":  userID,
			})
		})

		authorized.POST("/new_payment", func(c *gin.Context) {
			// TODO: update database to permit unlimited char/VAR for transaction_hashes, prevent page from duplicating data on refresh
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

			var splitDisplayTransactionHashesStr string
			if len(displayTransactionHashes) < 1 {
				c.Redirect(http.StatusFound, "/auth/new_payment_get")
			} else if len(displayTransactionHashes) > 1 {
				splitDisplayTransactionHashesStr = strings.Join(splitDisplayTransactionHashes, "\n")
				c.Set("splitDisplayTransactionHashesString", splitDisplayTransactionHashesStr)
				// fmt.Printf("splitDisplayTransactionHashesString %s\n", c.Get("splitDisplayTransactionHashesString"))
				fmt.Printf("splitDisplayTransactionHashesStr %s\n", splitDisplayTransactionHashesStr)
			}

			db.AutoMigrate(&User{})
			db.Save(&User{})

			c.HTML(http.StatusOK, "new_payment.html", gin.H{
				"title":                               "New Payment",
				"splitDisplayTransactionHashesString": splitDisplayTransactionHashesStr,
			})
		})

		authorized.GET("/transaction_get", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "transaction.html", gin.H{
				"title": "Transaction History",
				"user":  userID,
			})
		})

		authorized.POST("/transaction_history", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "transaction.html", gin.H{
				"title": "Transaction History",
				"user":  userID,
			})
		})

		authorized.GET("/new_contract_get", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "new_contract.html", gin.H{
				"title": "New Contract",
				"user":  userID,
			})
		})

		authorized.POST("/new_contract", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "new_contract.html", gin.H{
				"title": "New Contract",
				"user":  userID,
			})
		})

		authorized.GET("/payment_schedule", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("userId")

			c.HTML(http.StatusOK, "payment_schedule.tmpl", gin.H{
				"title": "Payment Schedule",
				"user":  userID,
			})
		})

		authorized.GET("/logout", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Delete("userId")
			session.Clear()
			session.Save()
			c.Redirect(http.StatusFound, "/index")
		})
	}

	// Listen and server on 0.0.0.0:8081
	router.Run(":8081")
}
