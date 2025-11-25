package question

import (
	"context"
	"encoding/json"
	"net/http"

	db "gamecraft-backend/prisma/db"
)

/* ---------- RESPONSE STRUCT ---------- */

type GetAllQuestionsResponse struct {
	Message   string      `json:"message"`
	Status    bool        `json:"status"`
	Data      interface{} `json:"data"`      // frontend-friendly key
	Count     int         `json:"count"`     // useful for pagination on frontend
}


func GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Allow only GET method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(GetAllQuestionsResponse{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	// Connect DB
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetAllQuestionsResponse{
			Message: "Database connection failed: " + err.Error(),
			Status:  false,
		})
		return
	}
	defer client.Prisma.Disconnect()

	// Fetch all questions
	ctx := context.Background()
	questions, err := client.QuestionRecords.FindMany().Exec(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetAllQuestionsResponse{
			Message: "Failed to fetch questions: " + err.Error(),
			Status:  false,
		})
		return
	}

	// Send clean JSON object for frontend rendering
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetAllQuestionsResponse{
		Message: "All questions fetched successfully",
		Status:  true,
		Data:    questions,          // easy array for frontend
		Count:   len(questions),     // helpful for listing UI/pagination
	})
}



// Code BY Gkrrr