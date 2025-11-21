package auths

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	db "gamecraft-backend/prisma/db"

	"github.com/golang-jwt/jwt/v5"
)




func Login(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		w.Header().Set("Context-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status: false,
		})
		return
	}




	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "failed to connect to server", http.StatusInternalServerError)
		return
	}
	defer client.Prisma.Disconnect()



	var user LoginController

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "invalid request body",
			Status:  false,
		})
		return
	}

	fmt.Println(user)


	existing, err := client.User.FindFirst(
		db.User.Or(
			db.User.Email.Equals(user.Id),
			db.User.Username.Equals(user.Id),
		),
	).Exec(context.Background())



	if err != nil || existing == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Message: "invalid credentials",
			Status:  false,
		})
		return
	}



	if !CheckPassword(user.Password, existing.Password) {
		fmt.Println("hii")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Message: "invalid credentials",
			Status:  false,
		})
		return
	}



	secret := []byte(os.Getenv("JWT_SECRET")) // set in .env
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existing.ID,
		"email":   existing.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // expires in 24h
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}




	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		HttpOnly: true,               // JS can’t read it
		Secure:   false,              // set true if using HTTPS
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	responseUserData := ResponseUserData{
		FirstName: existing.FirstName,
		LastName: existing.LastName,
		Email: existing.Email,
		Username : existing.Username,
	}

	fmt.Println(responseUserData)

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "login successful",
		Status:  true,
		Data: map[string]interface{}{
			"token": tokenString,
			"user":  responseUserData,
		},
	})
}