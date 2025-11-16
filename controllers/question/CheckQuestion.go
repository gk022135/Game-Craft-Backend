package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gamecraft-backend/controllers/helpers"
	"gamecraft-backend/middlewares"
	db "gamecraft-backend/prisma/db"

	"github.com/golang-jwt/jwt/v5"
)

func CkeckQustion(w http.ResponseWriter, r *http.Request) {
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


	userContainer, errUser := helpers.CreateMySQLContainer(existing.Username)
	testingContainer, errTest := helpers.CreateMySQLContainer(existing.Email)

	if errTest != nil || errUser != nil {
		fmt.Println("Error fetching questions: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message:  "failed to create contailer",
			Status:   false,
			TryLater: "please try again later",
		})
		return
	}
	defer helpers.DeleteContainer(userContainer)
	defer helpers.DeleteContainer(testingContainer)

	_, errDecUser := helpers.RunQuery(userContainer, question.StarterSchema)
	_, errDefUser := helpers.RunQuery(userContainer, question.StarterData)

	if errDecUser != nil || errDefUser != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "something error on checking",
			Status:  false,
		})
		return
	}

	_, errDecTest := helpers.RunQuery(testingContainer, question.StarterSchema)
	_, errDefTest := helpers.RunQuery(testingContainer, question.StarterData)

	if errDecTest != nil || errDefTest != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "something error on checking",
			Status:  false,
		})
		return
	}

	responseUser, errCheckUser := helpers.RunQuery(userContainer, solution.AnswerQuery)
	responseTest, errCheckTest := helpers.RunQuery(testingContainer, question.CorrectQuery)

	if errCheckTest != nil || errCheckUser != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "something error on checking",
			Status:  false,
		})
		return
	}



	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "questioned fetched successfully",
		Status:  true,
		Data:    result,
	})

}

