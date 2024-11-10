package middlewares

import (
	"fmt"
	"location-service/bin/config"
	"location-service/bin/pkg/utils"

	"github.com/labstack/echo/v4"
)

func VerifyBasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()
		if !ok {
			return utils.ResponseError(fmt.Errorf("Invalid username or password"), c)
		}
		if username == config.GetConfig().BasicAuthUsername && password == config.GetConfig().BasicAuthPassword {
			return next(c)
		}
		return utils.ResponseError(fmt.Errorf("Invalid username or password"), c)
	}
}
