package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// struttura utente

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	RoleID        string `json:"role_id"`
	AccountStatus string `json:"account_status"`
}

// FUNZIONI CRUD USERS

// CREATE USER

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var user User

	// decode JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Errore lettura dati utente", http.StatusBadRequest)
		return
	}

	// hash password bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(w, "Errore hash password", http.StatusInternalServerError)
		return
	}

	// query insert utente
	query := `
		INSERT INTO users (
			username,
			email,
			password_hash,
			first_name,
			last_name,
			birth_date,
			role_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = DB.Exec(
		query,
		user.Username,
		user.Email,
		string(hashedPassword),
		user.FirstName,
		user.LastName,
		user.BirthDate,
		user.RoleID,
	)

	if err != nil {
		http.Error(w, "Errore creazione utente", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Utente creato correttamente"))
}

// READ USERS

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo GET
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	// query database
	rows, err := DB.Query(`
		SELECT
			id,
			username,
			email,
			first_name,
			last_name,
			birth_date,
			role_id,
			account_status
		FROM users
		WHERE deleted_at IS NULL
	`)

	if err != nil {
		http.Error(w, "Errore recupero utenti", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	// loop righe database
	for rows.Next() {

		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.BirthDate,
			&user.RoleID,
			&user.AccountStatus,
		)

		if err != nil {
			fmt.Println("Errore scansione utente:", err)
			continue
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UPDATE USER

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var user User

	// decode JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Errore lettura dati utente", http.StatusBadRequest)
		return
	}

	// query update
	query := `
		UPDATE users
		SET
			email = $1,
			first_name = $2,
			last_name = $3,
			role_id = $4,
			account_status = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	_, err = DB.Exec(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.RoleID,
		user.AccountStatus,
		user.ID,
	)

	if err != nil {
		http.Error(w, "Errore aggiornamento utente", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Utente aggiornato correttamente"))
}

// DELETE USER-soft delete

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var user User

	// decode JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Errore lettura ID utente", http.StatusBadRequest)
		return
	}

	// soft delete account utente
	query := `
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err = DB.Exec(query, user.ID)

	if err != nil {
		http.Error(w, "Errore blocco utente", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Utente bloccato correttamente"))
}
