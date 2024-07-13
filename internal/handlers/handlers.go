package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"fourthtask/internal/db"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var regUser User
	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if regUser.Username == "" || regUser.Email == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`INSERT INTO users (username, email) VALUES ($1, $2)`, regUser.Username, regUser.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Failed to insert user:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandlerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.DB.Query(`SELECT username, email FROM users`)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Failed to query users:", err)
		return
	}
	defer rows.Close()

	var listUsers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username, &user.Email); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Failed to scan user:", err)
			return
		}
		listUsers = append(listUsers, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listUsers)
}
