package apiserver

import (
	"encoding/base64"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	v1 "go_learn/internal/biz/v1"
	"go_learn/internal/pkg/middleware"
	"go_learn/internal/pkg/middleware/auth"
	"go_learn/pkg/log"
	"net/http"
	"strings"
	"time"
)

const (
	Audience = "127.0.0.1"
	Issuer   = "apiserver"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {
		user, err := v1.NewUser().Get(username)
		if err != nil {
			return false
		}

		if err := user.Compare(password); err != nil {
			return false
		}

		return true
	})
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(auth.BasicStrategy), newJWTAuth().(auth.JWTStrategy))
}

func newJWTAuth() middleware.AuthStrategy {
	gJWT, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "jwt",
		SigningAlgorithm: "HS256",
		Key:              []byte("123456"),
		Timeout:          24 * time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      middleware.UsernameKey,
		// IdentityHandler IdentityKey 键的值由 IdentityHandler 函数返回
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return claims[jwt.IdentityKey]
		},
		// PayloadFunc 设置 JWT Token 中 Payload 部分的 iss、aud、sub、identity 字段 (添加额外业务相关的信息)
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			claims := jwt.MapClaims{
				"iss": Issuer,
				"aud": Audience,
			}
			if u, ok := data.(*v1.User); ok {
				claims[jwt.IdentityKey] = u.UserName
				claims["sub"] = u.UserName
			}
			return claims
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login loginInfo
			var err error

			if c.Request.Header.Get(auth.Authorization) != "" {
				login, err = parseWithHeader(c)
			} else {
				login, err = parseWithBody(c)
			}
			if err != nil {
				return "", jwt.ErrFailedAuthentication
			}

			user, err := v1.NewUser().Get(login.Username)
			if err != nil {
				log.Errorf("get user information failed: %s", err.Error())

				return "", jwt.ErrFailedAuthentication
			}

			if err := user.Compare(login.Password); err != nil {
				return "", jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "Login Success.",
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "Logout Success.",
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if u, ok := data.(string); ok && u == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: auth.Bearer,
		SendCookie:    true,
		TimeFunc:      time.Now,
	})

	return auth.NewJWTStrategy(*gJWT)
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	authRes := strings.SplitN(c.Request.Header.Get(auth.Authorization), " ", auth.AuthHeaderCount)
	if len(authRes) != auth.AuthHeaderCount || authRes[0] != auth.Basic {
		log.Errorf("get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(authRes[1])
	if err != nil {
		log.Errorf("decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		log.Errorf("parse payload failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	return loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Errorf("parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	return login, nil
}
