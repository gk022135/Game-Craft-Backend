// -------- Here I am Going to write the controller for contribut question --------//
package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db "gamecraft-backend/prisma/db"
)

/* ---------- REQUEST BODY STRUCT ---------- */

type ContributQuestionPayload struct {
	ContributedBy   string   `json:"contributed_by"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Topics          []string `json:"topics"`
	UsedTables      []string `json:"used_tables"`
	DifficultyLevel string   `json:"difficulty_level"`
	Rewards         *int     `json:"rewards"`
	Answer          *string  `json:"answer"`
	// fields kept exactly as per schema
}

func ContributeQuestion(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost { // FIXED: http.MethodePost → http.MethodPost
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	var payload ContributQuestionPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil { // FIXED syntax
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Invalid request body",
			Status:  false,
		})
		return
	}

	//-------extracting values from payload-------//
	ContributedBy := payload.ContributedBy
	Title := payload.Title
	Description := payload.Description
	Topics := payload.Topics
	UsedTables := payload.UsedTables
	DifficultyLevel := payload.DifficultyLevel

	// FIXED optional fields handling
	Rewards := 0
	if payload.Rewards != nil {
		Rewards = *payload.Rewards
	}

	Answer := "Not Provided"
	if payload.Answer != nil {
		Answer = *payload.Answer
	}

	//-------Connecting to database-------//
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil { // FIXED syntax
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Database connection failed: " + err.Error(),
			Status:  false,
		})
		return
	}
	defer client.Prisma.Disconnect()
	//-------Inserting contribut question into database-------//
	ctx := context.Background()
	_, err := client.QuestionRecords.CreateOne(
		db.QuestionRecords.ContributedBy.Set(ContributedBy),
		db.QuestionRecords.Title.Set(Title),
		db.QuestionRecords.Description.Set(Description),
		db.QuestionRecords.DifficultyLevel.Set(DifficultyLevel),
		db.QuestionRecords.UsedTables.Set(UsedTables),
		db.QuestionRecords.Topics.Set(Topics),
		db.QuestionRecords.Rewards.Set(Rewards),
		db.QuestionRecords.Answer.Set(Answer),
		db.QuestionRecords.TitleLowerCase.Set(strings.ToLower(Title)),
	).Exec(ctx) // FIXED missing Exec()

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Question contribution failed: " + err.Error(),
			Status:  false,
		})
		return
	}

	// all questions
	allQuestions, _ := client.QuestionRecords.FindMany().Exec(ctx) // FIXED: Findmany → FindMany

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Question contributed successfully",
		"status":    true,
		"questions": allQuestions,
	})
}

// Code By Gkrrr
