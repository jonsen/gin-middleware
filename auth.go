package middleware

import (
	"github.com/gin-gonic/gin"
)

// setup not need check authorize URL
//
//  	urls := []string{"/login", "/users", "/firewall"}
//	    router.Use(middleware.AuthRequest(urls...))
//
func AuthRequest(disUrls ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			match bool
		)
		ctx := Default(c)

		if u, has := ctx.SessGetValue("AuthInfo"); has && u != nil {
			ctx.Set("AuthData", u)
			c.Next()
			return
		}

		urlLen := len(disUrls)

		if urlLen < 1 {
			disUrls[0] = "/login"
			urlLen = 1
		}

		for _, url := range disUrls {
			if c.Request.RequestURI == url {
				match = true
				break
			}
		}

		if !match && c.Request.RequestURI != "/login" {
			c.Abort()
			c.Redirect(302, disUrls[0])
		}

	}
}
