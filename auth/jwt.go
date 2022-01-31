package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var JWTSecret = []byte("@#^T@#YHEWSE$Y&^#H")

func GenerateJWT(id uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}

type whiteList struct {
	path   string
	method string
}

var authWhiteList []whiteList

func AddToWhiteList(path string, method string) {
	if authWhiteList == nil {
		authWhiteList = make([]whiteList, 0)
	}
	authWhiteList = append(authWhiteList, whiteList{path, method})
}

func skipper(c echo.Context) bool {
	for _, v := range authWhiteList {
		if c.Path() == v.path && c.Request().Method == v.method {
			return true
		}
	}
	return false
}

func JWT() echo.MiddlewareFunc {
	c := middleware.DefaultJWTConfig
	c.SigningKey = JWTSecret
	c.Skipper = skipper
	return middleware.JWTWithConfig(c)
}
