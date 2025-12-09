package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "gamecraft-backend/prisma/db"
)

type TotalSolvedResponse struct {
	Email                string `json:"email"`
	TotalEasySolved      int    `json:"total_easy_solved"`
	TotalMediumSolved    int    `json:"total_medium_solved"`
	TotalHardSolved      int    `json:"total_hard_solved"`
	TotalAvailableEasy   int    `json:"total_available_easy"`
	TotalAvailableMedium int    `json:"total_available_medium"`
	TotalAvailableHard   int    `json:"total_available_hard"`
	PointsEarned         int    `json:"points_earned"`
}

type Response struct {
	Message  string      `json:"message"`
	Status   bool        `json:"status"`
	TryLater string      `json:"try_later,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func GetTotalSolved(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "Method not allowed",
		})
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "email is required",
		})
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "Failed to connect database",
		})
		return
	}
	defer client.Prisma.Disconnect()

	ctx := context.Background()

	// Fetch the user
	user, err := client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Status:  false,
			Message: "User not found",
		})
		return
	}

	//Total from QuestionsRecords
	var TotalAvailableEasy, TotalAvailableMedium, TotalAvailableHard int

	easyRecords, err := client.QuestionRecords.FindMany(
		db.QuestionRecords.DifficultyLevel.Equals("Easy"),
	).Exec(ctx)

	if err != nil {
		fmt.Println("Easy Count Error:", err)
	} else {
		TotalAvailableEasy = len(easyRecords)
	}

	mediumRecords, err := client.QuestionRecords.FindMany(
		db.QuestionRecords.DifficultyLevel.Equals("Medium"),
	).Exec(ctx)

	if err != nil {
		fmt.Println("Medium Count Error:", err)
	} else {
		TotalAvailableMedium = len(mediumRecords)
	}

	hardRecords, err := client.QuestionRecords.FindMany(
		db.QuestionRecords.DifficultyLevel.Equals("Hard"),
	).Exec(ctx)

	if err != nil {
		fmt.Println("Hard Count Error:", err)
	} else {
		TotalAvailableHard = len(hardRecords)
	}

	earnedPoints := user.EarnedPoints

	// Count the solved questions
	totalEasy := len(user.SolvedEasy)
	totalMedium := len(user.SolvedMedium)
	result := TotalSolvedResponse{
		Email:                user.Email,
		TotalEasySolved:      totalEasy,
		TotalMediumSolved:    totalMedium,
		TotalAvailableEasy:   TotalAvailableEasy,
		TotalAvailableMedium: TotalAvailableMedium,
		TotalAvailableHard:   TotalAvailableHard,
		PointsEarned:         earnedPoints,
	}
	

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  true,
		Message: "Total solved statistics fetched",
		Data:    result,
	})
}