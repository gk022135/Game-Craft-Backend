package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gamecraft-backend/controllers/helpers"
	"gamecraft-backend/middlewares"
	db "gamecraft-backend/prisma/db"

	"github.com/golang-jwt/jwt/v5"
)

// Response struct moved to QuestionStructures.go to avoid redeclaration

func RunQuestion(w http.ResponseWriter, r *http.Request) {
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
	
	// // Connect Prisma
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer client.Prisma.Disconnect()
	
	
	
	//  ////Get User
	claims, ok := r.Context().Value(middlewares.UserKey).(jwt.MapClaims)
    if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
	
	str := fmt.Sprintf("%v", claims["user_id"])
	num, _ := strconv.Atoi(str)
	existing, err := client.User.FindUnique(
		db.User.ID.Equals(num),
		).Exec(context.Background())
		
		
		
		if err != nil || existing == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Message: "invalid credentials",
				Status:  false,
			})
			return
	}
	
	
	//  ////Get Question
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

	question, err := client.Question.FindUnique(
		db.Question.ID.Equals(questionId),
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

	// Decode body → Go struct
	var solution SolutionController
	
	if err := json.NewDecoder(r.Body).Decode(&solution); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "invalid request body",
			Status:  false,
		})
		return
	}

	// Access fields
	query := solution.AnswerQuery
	
	// fmt.Println("User Email:", email)
	fmt.Println("Query:", query)
	
	// Execute queries using Prisma
	// Example: Loop and run each query

	if strings.TrimSpace(query) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Empty query",
			Status:  false,
			Result:  false,
		})
		return
	}
	
	
	
	testingResult, testingError := helpers.QueryRunner(question.StarterSchema, question.StarterData, question.CorrectQuery, question.EndingSchema)
	if testingError != nil {
		fmt.Println("testing error: ", testingError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "error in testing: " + testingResult,
			Data: testingError,
			Status:  false,
			Result: false,
		})
		return
	}
	
	fmt.Println("pahunch gaya 5")


	userResult, userError := helpers.QueryRunner(question.StarterSchema, question.StarterData, query, question.EndingSchema)
	if userError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(Response{
			Message: "Query is invalid or buggy",
			Data: userError,
			Result: false,
			Status:  false,
		})
		return
	}

	isCorrect := helpers.CompareResults(testingResult, userResult, false)

	fmt.Println("\n\n\n\n\n\n\n\n")
	fmt.Println(testingResult)
	fmt.Println("\n\n\n\n\n\n\n\n")
	fmt.Println(userResult)
	fmt.Println("\n\n\n\n\n\n\n\n")
	fmt.Println(isCorrect)
	fmt.Println("\n\n\n\n\n\n\n\n")

	var dt string

	if isCorrect {
		dt = "Correct query"
	} else {
		dt = "Incorrect query"
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Queries executed successfully",
		Status:  true,
		Result: isCorrect,
		Data: dt,
	})
}
