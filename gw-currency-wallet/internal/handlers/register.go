package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crypto/md5"
	"encoding/hex"
	. "gw-currency-wallet/internal/data"
	log "github.com/sirupsen/logrus"
)

func Register(ctx *gin.Context) {
	var user User
	if Err := ctx.ShouldBindJSON(&user); Err != nil {
		log.Error(Err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": Err.Error()})
		return
	}

	var usernameExist bool
	var emailExist bool
	Err = DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1);", user.Email).Scan(&emailExist)
	if Err != nil {
		log.Errorf("Error querying database for email existence: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	Err := DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1);", user.Username).Scan(&usernameExist)
	if Err != nil {
		log.Errorf("Error querying database for username existence: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	log.Info(usernameExist)

	if usernameExist || emailExist {
		log.Infof("Username: %s or email: %s already exist", user.Username, user.Email)
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	var walletId int
	Err = DB.QueryRow(ctx, "INSERT INTO wallets DEFAULT VALUES RETURNING id;").Scan(&walletId)
	if Err != nil {
		log.Errorf("Cannot create new wallet: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// тут можно было еще соль добавить
	hash := md5.New()
	hash.Write([]byte(user.Password))
	passHesh := hex.EncodeToString(hash.Sum(nil))

	_, Err = DB.Exec(ctx, "INSERT INTO users (username, password, wallet_id, email) VALUES ($1, $2, $3, $4);", user.Username, passHesh, walletId, user.Email)
	if Err != nil {
		log.Error("Insert user error: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
	return
}