package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

// refresh token
func GenerateRefreshToken(userID string, username string) (string, error) {

	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verifica del JWT
func VerifyJWT(tokenString string) (*jwt.Token, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	// prende header Authorization
	authHeader := r.Header.Get("Authorization")

	// verifica presenza token
	if authHeader == "" {
		http.Error(w, "Token mancante", http.StatusUnauthorized)
		return
	}

	// formato atteso:
	// Authorization: Bearer TOKEN
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// verifica JWT
	token, err := VerifyJWT(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Token non valido", http.StatusUnauthorized)
		return
	}

	// recupera claims JWT
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Errore claims JWT", http.StatusUnauthorized)
		return
	}

	// recupera user_id dal token
	userID, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Errore user_id JWT", http.StatusUnauthorized)
		return
	}

	// verifica ruolo admin
	isAdmin, err := CheckAdminRole(userID)
	if err != nil {
		http.Error(w, "Errore verifica ruolo", http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "Accesso negato: permessi insufficienti", http.StatusForbidden)
		return
	}

	// accesso consentito
	w.Write([]byte("Accesso autorizzato: ruolo admin verificato"))
}

// log degli accessi

func CreateAccessLog(userID string, accessResult string, ipAddress string) {

	_, err := DB.Exec(`
		INSERT INTO access_logs (
			user_id,
			access_result,
			ip_address
		)
		VALUES ($1, $2, $3)
	`, userID, accessResult, ipAddress)

	if err != nil {
		fmt.Println("Errore creazione access log:", err)
	}
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

		// access log per tentativo fallito

		CreateAccessLog(userID, "failed", r.RemoteAddr)

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

	// access log per login riuscito

	CreateAccessLog(userID, "success", r.RemoteAddr)

	// generazione access token
	accessToken, err := GenerateJWT(userID, loginData.Username)
	if err != nil {
		http.Error(w, "Errore generazione access token", http.StatusInternalServerError)
		return
	}

	// generazione refresh token
	refreshToken, err := GenerateRefreshToken(userID, loginData.Username)
	if err != nil {
		http.Error(w, "Errore generazione refresh token", http.StatusInternalServerError)
		return
	}

	// risposta JSON con token
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// funzione RBAC
func CheckAdminRole(userID string) (bool, error) {

	var roleName string

	query := `
		SELECT roles.role_name
		FROM users
		JOIN roles
			ON users.role_id = roles.id
		WHERE users.id = $1
	`

	err := DB.QueryRow(query, userID).Scan(&roleName)
	if err != nil {
		return false, err
	}

	return roleName == "admin", nil
}
