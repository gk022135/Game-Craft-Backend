package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "gamecraft-backend/prisma/db"
)

// Define request body structure
type QuestionRequest struct {
	Queries []string `json:"queries"`
	Email   string   `json:"email"`
}

type Person struct {
    PersonID  int     `json:"personid"`
    LastName  string  `json:"lastname"`
    FirstName string  `json:"firstname"`
    Address   string  `json:"address"`
    City      string  `json:"city"`
}

// Response struct moved to QuestionStructures.go to avoid redeclaration

func RunQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RunQuestion called")
	// Allow only POST
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	// Connect Prisma
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer client.Prisma.Disconnect()

	// Decode body → Go struct
	var questionBody QuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&questionBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "invalid request body",
			Status:  false,
		})
		return
	}

	// Access fields
	queries := questionBody.Queries
	email := questionBody.Email

	fmt.Println("User Email:", email)
	fmt.Println("Queries:", queries)

	// Execute queries using Prisma
	// Example: Loop and run each query
	results := []interface{}{}

	for _, q := range queries {
		people := []map[string]interface{}{}
		err := client.Prisma.QueryRaw(q).Exec(context.Background(), &people)
		if err != nil {
			json.NewEncoder(w).Encode(Response{
				Message: fmt.Sprintf("Query failed: %s", err.Error()),
				Status:  false,
			})
			return
		}
		fmt.Println("Query Results:", people)
		// Exec does not return query rows directly; record execution status
		results = append(results, map[string]interface{}{
			"query":  q,
			"ans":    people,
			"status": "executed",
		})
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Queries executed successfully",
		Data:    results,
		Status:  true,
	})
}
