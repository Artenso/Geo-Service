package token

import (
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func Init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWTSECRET")), nil)
}

func Generate(name string) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"username": name})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}
