package model

import (
	"log"
	"github.com/bndr/gopencils"
)

func LoginToEssentier(url, username, password string) string {
	essentierRest := gopencils.Api(url + "/essentier-rest")
	token := &JwtToken{}
	loginData := &LoginCredential{Email: username, Password: password}
	_, err := essentierRest.Res("login", token).Post(loginData)
	if err != nil {
		log.Printf("Failed to call the login rest api. Error is: %#v", err)
	}
	log.Printf("Received token is: %#v", token)
	return token.Token
}

type JwtToken struct {
	Token string `json:"token" form:"token"`
}

type LoginCredential struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
