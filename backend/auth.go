package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// funzione per la generazione del JWT
func GenerateJWT(userID string, username string) (string, error) {

	// carica il file .env
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	// legge la secret key dal file .env
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// accettazione esclusiva del metodo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var loginData LoginRequest

	// decoder del JSON inviato dal frontend
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Errore lettura dati login", http.StatusBadRequest)
		return
	}

	var userID string
	var passwordHash string
	var accountStatus string
	var failedAttempts int

	// ricerca dell'utente nel database
	query := `
		SELECT id, password_hash, account_status, failed_login_attempts
		FROM users
		WHERE username = $1
	`

	err = DB.QueryRow(query, loginData.Username).Scan(
		&userID,
		&passwordHash,
		&accountStatus,
		&failedAttempts,
	)

	// utente non trovato
	if err == sql.ErrNoRows {
		http.Error(w, "Utente non trovato", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Errore database", http.StatusInternalServerError)
		return
	}

	// account bloccato
	if accountStatus == "blocked" {
		http.Error(w, "Account bloccato", http.StatusForbidden)
		return
	}

	// verifica password con bcrypt
	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(loginData.Password),
	)

	if err != nil {

		// incrementa i tentativi falliti
		failedAttempts++

		// se arriva a 5 blocca l’account
		if failedAttempts >= 5 {
			_, err = DB.Exec(`
				UPDATE users
				SET failed_login_attempts = $1,
				    account_status = 'blocked'
				WHERE username = $2
			`, failedAttempts, loginData.Username)

			if err != nil {
				http.Error(w, "Errore aggiornamento account", http.StatusInternalServerError)
				return
			}

			http.Error(w, "Numero massimo di tentativi superato. Account bloccato", http.StatusForbidden)
			return
		}

		// aggiorna il numero di tentativi falliti
		_, err = DB.Exec(`
			UPDATE users
			SET failed_login_attempts = $1
			WHERE username = $2
		`, failedAttempts, loginData.Username)

		if err != nil {
			http.Error(w, "Errore aggiornamento tentativi", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Password non corretta", http.StatusUnauthorized)
		return
	}

	// reset tentativi falliti dopo login corretto
	_, err = DB.Exec(`
		UPDATE users
		SET failed_login_attempts = 0,
		    last_login = CURRENT_TIMESTAMP
		WHERE username = $1
	`, loginData.Username)

	if err != nil {
		http.Error(w, "Errore reset tentativi", http.StatusInternalServerError)
		return
	}

	// generazione JWT
	token, err := GenerateJWT(userID, loginData.Username)
	if err != nil {
		http.Error(w, "Errore generazione token", http.StatusInternalServerError)
		return
	}

	// risposta JSON con token
	response := map[string]string{
		"access_token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
