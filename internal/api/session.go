package api

import (
	"errors"
	"fmt"
	"wishmill/internal/logger"

	"github.com/golang-jwt/jwt"
)

func createToken(id int64, name string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"name":  name,
		"email": email,
	})

	tokenString, err := token.SignedString(token_secret)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	return tokenString, err
}

func verifyToken(tokenString string) (user *User, err error, tokenNotvalid bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return token_secret, nil
	})
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err, false
	}
	if !token.Valid {
		return nil, errors.New("token not valid"), true
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token not valid"), true
	}
	return &User{Id: int64(claims["id"].(float64)), Name: claims["name"].(string), Email: claims["email"].(string)}, nil, false
}
