package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "gamecraft-backend/prisma/db"
)

type QuestionResponseGk struct {
	Id           int      `json:"Id"`
	Title        string   `json:"Title"`
	Description  string   `json:"Description"`
	UsedTables   []string `json:"UsedTables"`
	ContributedBy string  `json:"ContributedBy,omitempty"`
	Points       int      `json:"Points,omitempty"`
	Answer 	 string   `json:"Answer,omitempty"`
}

func GetQustion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
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


	id := r.URL.Query().Get("id")
	questionId, err := strconv.Atoi(id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Invalid query parameters",
			Status: false,
		})
		return
	}

	question, err := client.QuestionRecords.FindUnique(
		db.QuestionRecords.ID.Equals(questionId),
	).Exec(context.Background())

	if err != nil {
		fmt.Println("Error fetching questions: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message:  "failed to fetch questions",
			Status:   false,
			TryLater: "please try again later",
		})
		return
	}

	// Map into a lightweight DTO
	var answer string
	if val, ok := question.Answer(); ok {
		answer = string(val)
	}

	result := QuestionResponseGk{
		Id:          question.ID,
		Title:       question.Title,
		Description: question.Description,
		UsedTables:  question.UsedTables,
		Points: func() int {
			if val, ok := question.Rewards(); ok {
				return int(val)
			}
			return 0
		}(),
		ContributedBy: question.ContributedBy,
		Answer:        answer, // optional
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "questioned fetched successfully",
		Status:  true,
		Data:    result,
	})

}

