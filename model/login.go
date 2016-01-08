package model

import "github.com/essentier/gopencils"

func LoginToEssentier(url, username, password string) (string, error) {
	essentierRest := gopencils.Api(url + "/essentier-rest")
	token := &JwtToken{}
	loginData := &LoginCredential{Email: username, Password: password}
	_, err := essentierRest.NewChildResource("login", token).Post(loginData)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

type JwtToken struct {
	Token string `json:"token" form:"token"`
}

type LoginCredential struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
