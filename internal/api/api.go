package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"wishmill/internal/config"
	"wishmill/internal/db"
	"wishmill/internal/logger"

	docs "wishmill/docs"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/oauth2"
)

// @title Wishmill API
// @version 1.0
// @BasePath /_wishmill/v1

var providers map[string]Oidc_provider

var token_secret []byte

type Errormsg struct {
	Message string `json:"message"`
}

type Oidc_provider struct {
	provider       *oidc.Provider
	config         oauth2.Config
	username_claim string
}

type User struct {
	Name         string `db:"name" json:"name"`
	Email        string `db:"email" json:"email"`
	AuthProvider string `db:"auth_provider" json:"-"`
	Sub          string `db:"sub" json:"-"`
	Id           int64  `db:"id" json:"-"`
}

type LoginBody struct {
	Provider    string `json:"provider" binding:"required"`
	Code        string `json:"code" binding:"required"`
	RedirectURL string `json:"redirect_url" binding:"required"`
}

type Session struct {
	Token string `json:"token" binding:"required"`
}

type Claim struct {
	Email              string `json:"email"`
	Sub                string `json:"sub"`
	Name               string `json:"name"`
	Issuer             string `json:"iss"`
	Preferred_Username string `json:"preferred_username"`
	Username           string `json:"username"`
}

func Init() {
	logger.DebugLogger.Println("api: init: Initializing api")
	var err error
	//Initialize o1c providers
	providers = map[string]Oidc_provider{}
	for _, a := range config.Config.Oidc_provider {
		o := Oidc_provider{}
		o.provider, err = oidc.NewProvider(context.TODO(), a.Url)
		if err != nil {
			logger.FatalLogger.Panicln("Failed to initialize oidc providers")
		}
		o.config = oauth2.Config{
			ClientID:     a.Client_id,
			ClientSecret: a.Client_secret,
			Endpoint:     o.provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		}
		o.username_claim = a.Username_claim
		providers[a.Name] = o
	}
	//Initialize Token Secret

	//Get token from yaml

	if config.Config.Token_secret != "" {
		logger.DebugLogger.Println("api: init: Getting token_secret from yaml")
		token_secret, err = base64.StdEncoding.DecodeString(config.Config.Token_secret)
		if err != nil {
			logger.FatalLogger.Println("Could not decode token_secret from yaml config")
			logger.FatalLogger.Panicln(err)
		}
		if len(token_secret) < 32 {
			logger.FatalLogger.Panicln("api: init: Token_secret from yaml config is not at least 256 bits long")
		}
	}

	//Get Token from DB
	logger.DebugLogger.Println("api: init: Getting token_secret from db")
	if config.Config.Token_secret == "" {
		rows, err := db.Db.Query("SELECT value FROM config WHERE key='token_secret'")
		if err != nil {
			logger.FatalLogger.Panicln(err)
		}
		defer rows.Close()
		next := rows.Next()
		//Token already in db
		if next {
			logger.DebugLogger.Println("api: init: token_secret already in DB, loading token_secret from db")
			var token_secret_string string
			rows.Scan(&token_secret_string)
			if err != nil {
				logger.FatalLogger.Panicln(err)
			}
			token_secret, err = base64.StdEncoding.DecodeString(token_secret_string)
			if err != nil {
				logger.FatalLogger.Panicln(err)
			}
		}
		//Token not in DB yet
		if !next {
			logger.DebugLogger.Println("api: init: No token in DB yet, generating token_secret")
			token_secret = make([]byte, 64)
			_, err = rand.Read(token_secret)
			if err != nil {
				logger.FatalLogger.Panicln(err)
			}
			logger.DebugLogger.Println("api: init: Writing token_secret to db")
			_, err = db.Db.Exec("INSERT INTO config(key, value) VALUES ('token_secret', $1)", base64.StdEncoding.EncodeToString(token_secret))
			if err != nil {
				logger.FatalLogger.Panicln(err)
			}
		}

	}
	//Initialize gin
	logger.DebugLogger.Println("api: init: Initializing gin")
	if !config.Config.DevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/_wishmill/v1/health", getHealth)
	router.GET("/_wishmill/v1/auth/oidc_providers", getOidc_Providers)
	router.POST("/_wishmill/v1/auth/obtainToken", obtainToken)
	router.POST("/_wishmill/v1/auth/checkToken", checkToken)
	docs.SwaggerInfo.BasePath = "/_wishmill/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	go router.Run(":8080")
	logger.DebugLogger.Println("api: init: Finished initializing api")
}

// getHealth godoc
// @Summary Get app health
// @Despcription Checks db connection
// @Success 200
// @Failure 500
// @Router /health [get]
func getHealth(c *gin.Context) {
	logger.DebugLogger.Println("api: getHealth is called")
	health := db.GetDbHealth()
	if health {
		c.String(200, "HEALTHY")
		logger.DebugLogger.Println("api: getHealth is OK")
		return
	}
	if !health {
		c.String(500, "NOT HEALTHY")
		logger.WarningLogger.Println("api: getHealth is NOT OK")
		return
	}
}

// @Summary Get oidc authentication providers
// @Despcription Get a list of oidc authentication provides, that can be used for authentication
// @Success 200 {array} config.Oidc_provider
// @Router /auth/oidc_providers [get]
func getOidc_Providers(c *gin.Context) {
	logger.DebugLogger.Println("api: getOidc_Servers is called")
	c.IndentedJSON(200, config.Config.Oidc_provider)
	logger.DebugLogger.Println("api: getOidc_Servers finished without error")
}

// obtainToken godoc
// @Summary Generate a session token
// @Despcription Perform a login action and create a session
// @Success 200 {object} Session
// @Failure 400 {object} Errormsg
// @Failure 500 {object} Errormsg
// @Router /auth/obtainToken [post]
// @Param Login body LoginBody true "Authorize data from oidc provider"
func obtainToken(c *gin.Context) {
	logger.DebugLogger.Println("api: login is called")

	var body LoginBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.IndentedJSON(400, Errormsg{Message: err.Error()})
		return
	}

	o, ok := providers[body.Provider]
	if !ok {
		logger.WarningLogger.Println("api: login: Non existant provider was called")
		c.IndentedJSON(400, Errormsg{Message: "Provider does not exist"})
		return
	}
	o.config.RedirectURL = body.RedirectURL
	token, err := o.config.Exchange(context.TODO(), body.Code)
	if err != nil {
		logger.ErrorLogger.Println("api: login: Failed to obtain token")
		logger.ErrorLogger.Println(err)
		c.IndentedJSON(500, Errormsg{Message: "Failed to obtain token"})
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		logger.ErrorLogger.Println("api: login: id_token is missing")
		c.IndentedJSON(500, Errormsg{Message: "Failed to obtain token"})
		return
	}

	verifier := o.provider.Verifier(&oidc.Config{ClientID: o.config.ClientID})
	idToken, err := verifier.Verify(context.TODO(), rawIDToken)
	if err != nil {
		logger.ErrorLogger.Println("api: login: Failed to verify token")
		logger.ErrorLogger.Println(err)
		c.IndentedJSON(500, Errormsg{Message: "Failed to obtain token"})
		return
	}

	claims := Claim{}
	err = idToken.Claims(&claims)
	if err != nil {
		logger.ErrorLogger.Println("api: login: Failed to get claims")
		logger.ErrorLogger.Println(err)
		c.IndentedJSON(500, Errormsg{Message: "Failed to obtain token"})
		return
	}

	var name string

	if o.username_claim == "email" {
		name = claims.Email
	}
	if o.username_claim == "preferred_username" {
		name = claims.Preferred_Username
	}
	if o.username_claim == "username" {
		logger.DebugLogger.Println("api: Usering claim username: " + claims.Username)
		name = claims.Username
	}
	if o.username_claim == "sub" {
		name = claims.Sub
	}
	if o.username_claim == "name" {
		name = claims.Name
	}

	id, err := registerUser(claims.Sub, body.Provider, name, claims.Email)
	if err != nil {
		c.IndentedJSON(500, Errormsg{Message: "Failed to create session"})
		return
	}

	token2, err := createToken(id, claims.Name, claims.Email)
	if err != nil {
		c.IndentedJSON(500, Errormsg{Message: "Failed to create session"})
		return
	}

	c.IndentedJSON(200, Session{Token: token2})
}

// checkToken godoc
// @Summary Check token validity
// @Despcription Check token validity and get user info
// @Success 200 {object} User
// @Failure 400 {object} Errormsg
// @Failure 500 {object} Errormsg
// @Router /auth/checkToken [post]
// @Param Session body Session true "Token"
func checkToken(c *gin.Context) {
	var session Session
	err := c.ShouldBindJSON(&session)
	if err != nil {
		c.IndentedJSON(400, Errormsg{Message: err.Error()})
		return
	}

	user, err, invalid := verifyToken(session.Token)
	if invalid || (err != nil) {
		c.IndentedJSON(400, Errormsg{Message: "Token not valid"})
		return
	}
	c.IndentedJSON(200, user)
}
