package tokens

import (
	"4microservice/api_gateway/config"
	"4microservice/api_gateway/pkg/logger"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTHandler struct {
	Sub       string
	Exp       string
	Role      string
	Iat       string
	SignInKey string
	Log       *logger.Logger
	Token     string
	Timeout   int
	cfg       config.Config
}
type CustomClaims struct {
	*jwt.Token
	Role string
	Name string
	Sub  string
	Exp  float64
	Iat  float64
}

func (jwtHandler *JWTHandler) GenerateAuthJWT(role string, name string, sub string) (access, refresh string, err error) {

	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
		rtClaims     jwt.MapClaims
	)

	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)
	claims = accessToken.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["sub"] = sub
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(3600)).Unix()
	claims["iat"] = time.Now().Unix()

	access, err = accessToken.SignedString([]byte("nodirbek"))
	if err != nil {
		log.Fatal("error while generating access token   ", err)
		return
	}
	rtClaims = refreshToken.Claims.(jwt.MapClaims)
	rtClaims["name"] = name
	rtClaims["sub"] = sub
	rtClaims["role"] = role
	rtClaims["exp"] = time.Now().Add(time.Minute * time.Duration(3600)).Unix()
	rtClaims["iat"] = time.Now().Unix()

	refresh, err = refreshToken.SignedString([]byte("nodirbek"))
	if err != nil {
		log.Fatal("error while generating refresh token", err)
	}
	return access, refresh, nil
}

func (jwtHandler *JWTHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SignInKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		log.Fatal("invalid jwt token")
		return nil, err
	}
	return claims, nil
}
