package main

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// middleware autenticazione JWT

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// prende Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Token mancante", http.StatusUnauthorized)
			return
		}

		// rimuove Bearer
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// verifica token
		token, err := VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Token non valido", http.StatusUnauthorized)
			return
		}

		// verifica claims JWT
		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Errore claims JWT", http.StatusUnauthorized)
			return
		}

		// passa al prossimo handler
		next(w, r)
	}
}
