package handlers

import (
	"fmt"
	"net/http"
	"task-golang-db/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthInterface interface {
	Login(*gin.Context)
	Upsert(*gin.Context)
}

type authImplement struct {
	db         *gorm.DB
	signingKey []byte
}

func NewAuth(db *gorm.DB, signingKey []byte) AuthInterface {
	return &authImplement{
		db,
		signingKey,
	}
}

type authLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *authImplement) Login(c *gin.Context) {
	payload := authLoginPayload{}

	// parsing JSON payload to struct models
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Validate username to get auth data
	auth := models.Auth{}
	if err := a.db.Where("username = ?",
		payload.Username).
		First(&auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "username salah",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(payload.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "password salah",
		})
		return
	}

	// Login is valid
	token, err := a.createJWT(&auth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode) // Set SameSite attribute (for cross-origin requests)
	// c.SetCookie("auth_token", token, 3600*72, "/", "", false, true)

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%v Login Sukses", payload.Username),
		"data":    token,
	})
}

type authUpsertPayload struct {
	AccountID int64  `json:"account_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (a *authImplement) Upsert(c *gin.Context) {
	payload := authUpsertPayload{}

	// parsing JSON payload to struct models
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Hash Given Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Check AccountID is valid
	var account models.Account
	if err := a.db.First(&account, payload.AccountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Account Not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Prepare new auth data with new password
	auth := models.Auth{
		AccountID: payload.AccountID,
		Username:  payload.Username,
		Password:  string(hashed),
	}

	// Upsert auth data (Insert or Update if already exists)
	result := a.db.Clauses(
		clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"username", "password"}),
			Columns:   []clause.Column{{Name: "account_id"}},
		}).Create(&auth)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Create success",
		"data":    payload.Username,
	})
}

func (a *authImplement) createJWT(auth *models.Auth) (string, error) {
	// Create the jwt token signer
	token := jwt.New(jwt.SigningMethodHS256)

	// Add claims data or additional data (avoid to put secret information in the payload or header elements)
	claims := token.Claims.(jwt.MapClaims)
	claims["auth_id"] = auth.AuthID
	claims["account_id"] = auth.AccountID
	claims["username"] = auth.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expires in 72 hours

	// Encode
	tokenString, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}

	// Return the token
	return tokenString, nil
}
