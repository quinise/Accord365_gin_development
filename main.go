package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	//"io/ioutil"
	//"log"
	//"html/template"
	"net/http"

	//"os"
	"path/filepath"
	//"strings"
	"time"

	// "github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"

	// "github.com/accord365/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"

	// "github.com/markbates/goth/gothic"
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
// 		RedirectURL:  "http://127.0.0.1:8080/auth",
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
// 		},
// 		Endpoint: google.Endpoint,
// 	}
// }

// func indexHandler(c *gin.Context) {
// 	c.HTML(http.StatusOK, "index.html", gin.H{})
// }

// func getLoginURL(state string) string {
// 	return conf.AuthCodeURL(state)
// }

// // Google AuthHandler handles authentication of a user and initiates a session.
// func authHandler(c *gin.Context) {
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
// 		fmt.Println("This is where we're failing")
// 		return
// 	}
// 	defer userinfo.Body.Close()
// 	data, _ := ioutil.ReadAll(userinfo.Body)
// 	u := SessionUser{}
// 	if err = json.Unmarshal(data, &u); err != nil {
// 		log.Println(err)
// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
// 		return
// 	}

// Save unique id / email to session
// 	session.Set("userId", u.Sub)
// 	session.Set("userGivenName", u.GivenName)
// 	session.Set("userFamilyName", u.FamilyName)
// 	err = session.Save()

// 	userID := session.Get("userId")
// 	userGivenName := session.Get("userGivenName")
// 	userFamilyName := session.Get("userFamilyName")

// 	// query database for u.Email/userID
// 	idResult := User{}
// 	emptyUser := User{}
// 	// ToDo (production): instead of email call it identification string db field must be unique
// 	db.Where("email = ?", userID).First(&idResult)
// 	if idResult == emptyUser {
// 		// fmt.Println("idResult ", idResult)
// 		user := User{Name: userGivenName.(string) + " " + userFamilyName.(string), Provider: "google", ProviderID: "3", Email: userID.(string)}
// 		result := db.Create(&user)

// 		fmt.Println("create user Result ", result)
// 	}

// 	if err != nil {
// 		log.Println(err)
// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
// 		return
// 	}
// 	// var seen = false
// 	if session.Get("userId") == u.Sub {
// 		// seen := true
// 		log.Println("seen!")
// 	} else {
// 		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving user. Please try again."})
// 		return
// 	}

// 	// seen := false
// 	dsn := "root:obatala88@tcp(127.0.0.1:3306)/accord365_development?charset=utf8&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		// seen = true
// 		panic(err.Error())
// 	} else {
// 		results := db.Save(&u)
// 		if results != nil {
// 			log.Println(results)
// 			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving user. Please try again."})
// 			return
// 		}
// 	}

// 	c.Redirect(http.StatusFound, "/auth/dashboard")

// 	// c.HTML(http.StatusOK, "dashboard.html", gin.H{
// 	// 	"title":  "Dashboard",
// 	// 	"userId": session.Get("userId"),
// 	// 	"seen":   seen,
// 	// })
// }

// func loginHandler(c *gin.Context) {
// 	state = randToken()
// 	session := sessions.Default(c)
// 	session.Set("state", state)
// 	session.Save()
// 	c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
// }

func main() {

	key := "h)xx5i%Ob[E+K_b" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30     // 30 days
	isProd := false          // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	// gothic.Store = store

	googleProvider := google.New("649290674745-7u4vagpfm303plk8qo8sktpdb5vscl0l.apps.googleusercontent.com", "ChUeGtuJ1IrlIWeddHF3XS9g", "http://localhost:8080/auth/google/callback")
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
