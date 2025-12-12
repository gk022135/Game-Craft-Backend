package users

import (
	"context"
	"encoding/json"
	"net/http"

	db "gamecraft-backend/prisma/db"
)

/** ---------- REQUEST PAYLOAD ---------- */
type UpdateUserProfilePayload struct {
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	GithubUrl   string   `json:"githubUrl"`
	LinkedinUrl string   `json:"linkedinUrl"`
	WebsiteUrl  string   `json:"websiteUrl"`
	Skills      []string `json:"skills"`
	PhoneNumber string   `json:"phoneNumber"`
}

/** ---------- RESPONSE STRUCT ---------- */

/** ---------- HANDLER ---------- */
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "Only POST method allowed",
		})
		return
	}

	var payload UpdateUserProfilePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "Invalid JSON body",
		})
		return
	}

	if payload.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "Email is required",
		})
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect DB", http.StatusInternalServerError)
		return
	}
	defer client.Prisma.Disconnect()

	// -------- FETCH USER --------
	user, err := client.User.FindUnique(
		db.User.Email.Equals(payload.Email),
	).Exec(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "User not found",
		})
		return
	}

	// Prepare update values
	updateData := []db.UserSetParam{}

	if payload.Username != "" {
		updateData = append(updateData, db.User.Username.Set(payload.Username))
	}

	updateData = append(updateData,
		db.User.Links.Set([]string{
			"Github:" + payload.GithubUrl,
			"Linkedin:" + payload.LinkedinUrl,
			"Website:" + payload.WebsiteUrl,
			"Skills:" + sliceToCSV(payload.Skills),
		}),
	)

	if payload.PhoneNumber != "" {
		updateData = append(updateData, db.User.PhoneNumber.Set(payload.PhoneNumber))
	}

	// -------- UPDATE USER --------
	_, err = client.User.FindUnique(
		db.User.Email.Equals(payload.Email),
	).Update(
		updateData...,
	).Exec(context.Background())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "Failed to update user",
		})
		return
	}

	// -------- SEND RESPONSE --------
	json.NewEncoder(w).Encode(Response{
		Status:  true,
		Message: "User profile updated successfully",
		Data: map[string]interface{}{
			"email":       user.Email,
			"username":    payload.Username,
			"githubUrl":   payload.GithubUrl,
			"linkedinUrl": payload.LinkedinUrl,
			"websiteUrl":  payload.WebsiteUrl,
			"skills":      payload.Skills,
			"phone":       payload.PhoneNumber,
		},
	})
}

/** Converts []string → comma-separated string */
func sliceToCSV(arr []string) string {
	result := ""
	for i, v := range arr {
		if i > 0 {
			result += ","
		}
		result += v
	}
	return result
}
