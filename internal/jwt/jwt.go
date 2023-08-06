package jwt

import (
	"context"
	"errors"

	"github.com/go-chi/jwtauth/v5"
)

var (
	NoUserID = errors.New("no user id")
)

func GetTokenAuth(secret string) *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(secret), nil)
}

func MakeJWT(secret string, payload map[string]any) (string, error) {
	tokenAuth := GetTokenAuth(secret)
	_, tokenString, err := tokenAuth.Encode(payload)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func MakeJWTPayload(ID int64) map[string]any {
	return map[string]any{
		"ID": ID,
	}
}

func GetUserID(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	userIDInterface := claims["ID"]
	//JSON unmarshalls all number to float64
	if userID, ok := userIDInterface.(float64); ok {
		return int64(userID), nil
	}

	return 0, NoUserID
}
