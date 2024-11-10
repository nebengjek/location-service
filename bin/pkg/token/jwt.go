package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"location-service/bin/config"
	"location-service/bin/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

func Validate(ctx context.Context, publicKey string, tokenString string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			output <- utils.Result{Error: "Failed to parse public key"}
			return
		}
		audience := config.GetConfig().JwtAudience
		issuer := config.GetConfig().JwtIssuer
		algorithm := config.GetConfig().JwtAlgorithm
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{audience},
			Issuer:   issuer,
		}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok || token.Header["alg"] != algorithm {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				output <- utils.Result{Error: "token has been expired"}
				return
			}
			output <- utils.Result{Error: "token parsing error"}
			return
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			var tokenClaim Claim
			jsonString, _ := json.Marshal(claims)
			json.Unmarshal(jsonString, &tokenClaim)
			output <- utils.Result{Data: tokenClaim}
		} else {
			output <- utils.Result{Error: "Token is not valid!"}
			return
		}
	}()

	return output
}
