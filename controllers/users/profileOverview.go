package users

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"strconv"

	db "gamecraft-backend/prisma/db"
)

type FrontendUserResponse struct {
	FirstName    string   `json:"FirstName"`
	LastName     string   `json:"LastName"`
	Username     string   `json:"Username"`
	Email        string   `json:"email"`
	CurrentLevel int      `json:"currentLevel"`
	TotalPoints  int      `json:"totalPoints"`
	Badges       []string `json:"badges"`
	GamesPlayed  int      `json:"gamesPlayed"`
	GithubUrl    string   `json:"githubUrl"`
	LinkedinUrl  string   `json:"linkedinUrl"`
	WebsiteUrl   string   `json:"websiteUrl"`
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "Missing email parameter",
		})
		return
	}

	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "DB connection error",
		})
		return
	}
	defer client.Prisma.Disconnect()

	ctx := context.Background()

	user, err := client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil || user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "User not found",
		})
		return
	}

	// Extract values from Links[]
	var github, linkedin, website string
	gamesPlayed := 0
	var badges []string

	for _, item := range user.Links {
		if strings.HasPrefix(item, "github:") {
			github = strings.TrimPrefix(item, "github:")
		}
		if strings.HasPrefix(item, "linkedin:") {
			linkedin = strings.TrimPrefix(item, "linkedin:")
		}
		if strings.HasPrefix(item, "website:") {
			website = strings.TrimPrefix(item, "website:")
		}
		if strings.HasPrefix(item, "games_played:") {
			val := strings.TrimPrefix(item, "games_played:")
			g, _ := strconv.Atoi(val)
			gamesPlayed = g
		}
		if strings.HasPrefix(item, "badge:") {
			b := strings.TrimPrefix(item, "badge:")
			badges = append(badges, b)
		}
	}

	// Prepare frontend response
	resp := FrontendUserResponse{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		Email:        user.Email,
		CurrentLevel: user.CurrentLevel,
		TotalPoints:  user.EarnedPoints,
		Badges:       badges,
		GamesPlayed:  gamesPlayed,
		GithubUrl:    github,
		LinkedinUrl:  linkedin,
		WebsiteUrl:   website,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data":   resp,
	})
}
