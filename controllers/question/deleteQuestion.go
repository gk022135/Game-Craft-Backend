package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "gamecraft-backend/prisma/db"
)

func DeleteQuestionFromRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete { // FIXED: http.MethodePost → http.MethodPost
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
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

	_, err = client.Question.FindUnique(
		db.Question.ID.Equals(questionId),
	).Delete().Exec(context.Background())

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Question deleted successfully",
		Status:  true,
	})

}