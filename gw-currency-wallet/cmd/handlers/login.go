package handlers

import (
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"crypto/md5"
	"encoding/hex"
	. "gw-currency-wallet/cmd/data"
	log "github.com/sirupsen/logrus"
)

func Login(ctx *gin.Context) {
	var user User
	if Err := ctx.ShouldBindJSON(&user); Err != nil {
		log.Errorf("Field to bind JSON: %v", Err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": Err.Error()})
		return
	}

	var password string
	var userId int
	Err := DB.QueryRow(ctx, "SELECT password, id FROM users WHERE username = $1;", user.Username).Scan(&password, &userId)
	
	if Err != nil {
		log.Errorf("Cannot find username: %v", Err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	hash := md5.New()
	hash.Write([]byte(user.Password))
	passHesh := hex.EncodeToString(hash.Sum(nil))

	if passHesh != password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, Err := token.SignedString(SecretKey)
	if Err != nil {
		log.Errorf("Error to sign token: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": Err.Error()})
        return
	}

	_, Err = DB.Exec(ctx, "INSERT INTO JWTTokens (token, expiration, user_id) VALUES ($1, $2, $3)", tokenString, expirationTime, userId)
	if Err != nil {
		log.Errorf("Cannot add JWT token into database: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}