package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"location-service/bin/config"
	"location-service/bin/pkg/token"
	"location-service/bin/pkg/utils"

	"github.com/labstack/echo/v4"
)

func decodeKey(secret string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func VerifyBearer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := strings.TrimPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")

		if len(tokenString) == 0 {
			return utils.Response(nil, "Invalid token!", http.StatusUnauthorized, c)
		}

		publicKey, err := decodeKey(config.GetConfig().PublicKey)
		if err != nil {
			return utils.Response(nil, utils.ConvertString(err), http.StatusUnauthorized, c)
		}
		parsedToken := <-token.Validate(c.Request().Context(), publicKey, tokenString)
		if parsedToken.Error != nil {
			return utils.Response(nil, utils.ConvertString(parsedToken.Error), http.StatusUnauthorized, c)
		}
		data, _ := json.Marshal(parsedToken.Data)
		jsonData := []byte(data)
		var claim token.Claim
		json.Unmarshal(jsonData, &claim)
		c.Set("userId", claim.Sub)
		return next(c)
	}
}
