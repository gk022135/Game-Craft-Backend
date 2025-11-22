package question


import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "gamecraft-backend/prisma/db"
)

func SaveQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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



	var question QuestionController

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "invalid request body",
			Status:  false,
		})
		return
	}
	// fmt.Println("Question: ", question)





	newQuestion, err := client.Question.CreateOne(
		db.Question.Title.Set(question.Title),
		db.Question.Description.Set(question.Description),
		db.Question.StarterSchema.Set(question.StarterSchema),
		db.Question.StarterData.Set(question.StarterData),
		db.Question.CorrectQuery.Set(question.CorrectQuery),
		db.Question.EndingSchema.Set(question.EndingSchema),
	).Exec(context.Background())

	if err != nil {
		fmt.Println("Error saving question:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message:  "failed to save question",
			Status:   false,
			TryLater: "please try again later",
		})
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "question saved successfully",
		Status:  true,
		Data:    newQuestion,
	})

}

