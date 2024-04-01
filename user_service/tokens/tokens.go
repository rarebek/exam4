package tokens

import (
	"4microservice/user-service/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type JWTHandler struct {
	Sub       string
	Exp       string
	Iat       string
	SignInKey string
	Log       *logger.Logger
	Token     string
	Timeout   int
}
type CustomClaims struct {
	*jwt.Token
	Role string
	Name string
	Sub  string
	Exp  float64
	Iat  float64
}

func (jwtHandler *JWTHandler) GenerateAuthJWT() (access, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
		rtClaims     jwt.MapClaims
	)

	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)
	claims = accessToken.Claims.(jwt.MapClaims)
	claims["name"] = "temp"
	claims["sub"] = jwtHandler.Sub
	claims["role"] = "user"
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(jwtHandler.Timeout)).Unix()
	claims["iat"] = time.Now().Unix()

	access, err = accessToken.SignedString([]byte(jwtHandler.SignInKey))
	if err != nil {
		log.Fatal("error while generating access token", err)
		return
		return
	}
	rtClaims = refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = jwtHandler.Sub
	refresh, err = refreshToken.SignedString([]byte(jwtHandler.SignInKey))
	if err != nil {
		log.Fatal("error while generating refresh token", err)
	}
	return
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
