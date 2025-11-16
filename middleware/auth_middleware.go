package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var userClaimsKey = "UsCLKey"
var jwtSecret = []byte("SecretKey")

type UserClaims struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}

func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	v, ok := ctx.Value(userClaimsKey).(UserClaims)
	return &v, ok
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authValue := r.Header.Get("Authorization")
		if !strings.HasPrefix(authValue, "Bearer ") {
			http.Error(w, "Потерянный токен авторизации", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authValue, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Не валидный токен авторизации", http.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Не валидный токен Claims", http.StatusUnauthorized)
		}
		userMap, ok := claims["user"].(map[string]interface{})
		if !ok {
			panic("invalid user data in token")
		}
		id, ok := userMap["id"].(string)
		if !ok {
			http.Error(w, "invalid id in token", http.StatusUnauthorized)
			return
		}
		username, ok := userMap["username"].(string)
		if !ok {
			http.Error(w, "invalid username in token", http.StatusUnauthorized)
			return
		}
		user := UserClaims{
			ID:       id,
			UserName: username,
		}
		ctx := context.WithValue(r.Context(), userClaimsKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
