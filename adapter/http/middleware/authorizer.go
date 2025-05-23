package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	AccountID string `json:"sub"`
	jwt.RegisteredClaims
}

func JwtAuthorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Formato de token inválido")
			return
		}

		tokenStr := tokenParts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		fmt.Printf("error: %v", err)
		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Token invalido: "+err.Error())
			return
		}
		// if claims.AccountID != "" {
		// 	exists, err := chamar service pra verificar se conta existe
		// 	if err != nil {
		// 		fmt.Printf("Error checking account existence: %v\n", err)
		// 		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		// 		return
		// 	}
		// 	if !exists {
		// 		respondWithError(w, http.StatusUnauthorized, "Account not found")
		// 		return
		// 	}
		// }
		//
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, message)))
}
