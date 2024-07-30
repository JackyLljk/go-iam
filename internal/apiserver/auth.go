package apiserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"j-iam/internal/apiserver/store"
	"j-iam/internal/pkg/middleware"
	"j-iam/internal/pkg/middleware/auth"
	"j-iam/pkg/log"

	v1 "j-iam/internal/pkg/model/apiserver/v1"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/spf13/viper"
)

const (
	// 定义 jwt audience 字段（签发和验证时，aud 字段必须是一致的！）
	apiServerAudience = "j-iam"

	// 定义 jwt issuer 字段，即 jwt token 的签发者
	apiServerIssuer = "j-apiserver"
)

type loginInfo struct {
	Username string `form:"username" json:"username" validator:"required,username"`
	Password string `form:"password" json:"password" validator:"required,password"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {
		// fetch user from database
		user, err := store.Client().Users().Get(context.TODO(), username, metav1.GetOptions{})
		if err != nil {
			return false
		}

		// Compare the login password with the user password.
		if err := user.Compare(password); err != nil {
			return false
		}

		user.LoginedAt = time.Now()
		_ = store.Client().Users().Update(context.TODO(), user, metav1.UpdateOptions{})

		return true
	})
}

// newJWTAuth 初始化 GinJWTMiddleware
func newJWTAuth() middleware.AuthStrategy {
	ginJWT, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout"),
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"),
		Authenticator:    authenticator(),
		LoginResponse:    loginResponse(),
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, nil)
		},
		RefreshResponse: refreshResponse(),
		PayloadFunc:     payloadFunc(),
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return claims[jwt.IdentityKey]
		},
		IdentityKey:  middleware.UsernameKey,
		Authorizator: authorizator(),
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
		// TODO: HTTPStatusMessageFunc:
	})

	return auth.NewJWTStrategy(*ginJWT)
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(auth.BasicStrategy), newJWTAuth().(auth.JWTStrategy))
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login loginInfo
		var err error

		// support header and body both
		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}
		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		// Get the user information by the login username.
		user, err := store.Client().Users().Get(c, login.Username, metav1.GetOptions{})
		if err != nil {
			log.Errorf("get user information failed: %s", err.Error())

			return "", jwt.ErrFailedAuthentication
		}

		// Compare the login password with the user password.
		if err := user.Compare(login.Password); err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		user.LoginedAt = time.Now()
		_ = store.Client().Users().Update(c, user, metav1.UpdateOptions{})

		// 登录认证通过后，将 username 添加到 gin.Context
		c.Set(middleware.UsernameKey, login.Username)

		return user, nil
	}
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(authorization) != 2 || authorization[0] != "Basic" {
		log.Errorf("get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(authorization[1])
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

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		//log.Info(c.GetString("username"))
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		fmt.Println(data)
		claims := jwt.MapClaims{
			"iss": apiServerIssuer,
			"aud": apiServerAudience,
		}
		if u, ok := data.(*v1.User); ok {
			claims[jwt.IdentityKey] = u.Name
			claims["sub"] = u.Name
			fmt.Println(claims)
		}

		return claims
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(string); ok {
			log.L(c).Infof("user `%s` is authenticated.", v)

			return true
		}

		return false
	}
}
