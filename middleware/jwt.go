package middleware

//
//import (
//	"context"
//	"github.com/cloudwego/hertz/pkg/app"
//	"github.com/hertz-contrib/jwt"
//	"time"
//)
//
//func JWT() {
//
//	jwt.New(&jwt.HertzJWTMiddleware{
//		Realm:      "test Realm",
//		Key:        []byte("secrect key"),
//		Timeout:    time.Hour,
//		MaxRefresh: time.Hour,
//		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
//
//		},
//		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
//			return false
//		},
//		PayloadFunc: func(data interface{}) jwt.MapClaims {
//
//		},
//		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
//
//		},
//		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
//
//		},
//		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
//
//		},
//		RefreshResponse: func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
//
//		},
//		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
//			return nil
//		},
//		TokenLookup: "",
//
//		WithoutDefaultTokenHeadName: false,
//		TimeFunc:                    nil,
//		HTTPStatusMessageFunc:       nil,
//		PrivKeyFile:                 "",
//		PrivKeyBytes:                nil,
//		PubKeyFile:                  "",
//		PrivateKeyPassphrase:        "",
//		PubKeyBytes:                 nil,
//		SendCookie:                  false,
//		CookieMaxAge:                0,
//		SecureCookie:                false,
//		CookieHTTPOnly:              false,
//		CookieDomain:                "",
//		SendAuthorization:           false,
//		DisabledAbort:               false,
//		CookieName:                  "",
//		CookieSameSite:              0,
//		ParseOptions:                nil,
//	})
//}
