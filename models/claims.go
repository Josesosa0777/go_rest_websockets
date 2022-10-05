package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId             string `json:"userId"` // el user será capaz de identicarse a través del user id que irá en un token
	jwt.StandardClaims        // al poner StandardClaims indico que AppClaims tiene todas las propiedades que están definidas en StandardClaims
}
