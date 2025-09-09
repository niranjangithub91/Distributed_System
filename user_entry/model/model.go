package model

import "github.com/dgrijalva/jwt-go"

type Data map[string]interface{}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Claims struct {
	Data string `json:"username"`
	jwt.StandardClaims
}
type Collection struct {
	Name      string `json:"username"`
	File_name string `json:"file_name"`
	Port      int    `json:"port"`
	Chunk     int    `json:"chunk"`
}
