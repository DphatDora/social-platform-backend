package util

import (
	"net/http"
	"social-platform-backend/config"

	"github.com/gin-gonic/gin"
)

func SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	conf := config.GetConfig()
	maxAge := conf.Auth.RefreshTokenExpirationDays * 24 * 60 * 60

	if conf.App.Debug {
		c.SetSameSite(http.SameSiteLaxMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}

	c.SetCookie(
		"refreshToken",          // name
		refreshToken,            // value
		maxAge,                  // max age in seconds
		"/",                     // path
		"",                      // domain (empty = current domain)
		conf.App.Debug == false, // secure (HTTPS only in production)
		true,                    // httpOnly
	)
}

func ClearRefreshTokenCookie(c *gin.Context) {
	conf := config.GetConfig()

	if conf.App.Debug {
		c.SetSameSite(http.SameSiteLaxMode)
	} else {
		c.SetSameSite(http.SameSiteNoneMode)
	}

	c.SetCookie(
		"refreshToken",
		"",
		-1, // MaxAge -1 deletes the cookie
		"/",
		"",
		conf.App.Debug == false,
		true,
	)
}
