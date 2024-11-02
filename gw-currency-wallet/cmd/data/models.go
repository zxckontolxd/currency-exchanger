package data

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
}

type Claims struct {
	Username string `json:"name"`
	jwt.StandardClaims
}

type Deposite struct {
	Amount float64 `json:"amount"`
	Currency string `json:"currency"` 
}

type Exchanger struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency string `json:"to_currency"`
	Amount float64 `json:"amount"`
}