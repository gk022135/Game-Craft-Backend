package users

import (
	"context"
	"encoding/json"
	"net/http"

	db "gamecraft-backend/prisma/db"
)

func GetUserActivity(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email parameter is required", http.StatusBadRequest)
		return
	}

	client := db.NewClient()
	if err := client.Connect(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect()

	ctx := context.Background()

	// Find the user
	user, err := client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch activity logs
	records, err := client.UserActivityLog.FindMany(
		db.UserActivityLog.UserID.Equals(user.ID),
	).Exec(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  Group by date yyyy-mm-dd
	activityMap := map[string]int{}
	for _, record := range records {
		day := record.Timestamp.Format("2006-01-02")
		activityMap[day]++
	}

	// Convert map → slice
	type Activity struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var result []Activity
	for date, count := range activityMap {
		result = append(result, Activity{Date: date, Count: count})
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
