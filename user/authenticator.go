package user

import (
	"crypto/rsa"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

type Authenticator struct {
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
}

func NewAuthenticator(privKeyPath, pubKeyPath string) *Authenticator {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	return &Authenticator{verifyKey, signKey}
}

func (a *Authenticator) CreateToken(username string) (token string, err error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims["AccessToken"] = "level1"
	t.Claims["CustomUserInfo"] = struct {
		Name string
		Kind string
	}{username, "human"}

	t.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	return t.SignedString(a.signKey)
}

func (a *Authenticator) Validate(token string) (err error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return a.verifyKey, nil
	})

	switch err.(type) {
	case nil: // no error

		if !t.Valid { // but may still be invalid
			return
		}

	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			return
		default:
			return
		}

	default: // something else went wrong
		return
	}

	return err
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
