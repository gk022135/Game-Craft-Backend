/*
User Payload Have
 1. question title and id and user answer query ----> from Tales_Info table in Main Database
 2. by question id and title we can find the question and used tables for tthis question
 3. also get the correct annswer from main database

 ----- Here all need of question done --------
===========================================
-------We Have to Run the user Query and Correct Answer Query-------

after runing the both query we have to compare the results
if both results are same then we have to return success response
if both results are different then we have to return failure response with correct answer

-----------------------------------------------------------
*/

package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	db "gamecraft-backend/prisma/db"
	db2 "gamecraft-backend/prisma_games/prisma_games_client"
)

/* ---------- REQUEST PAYLOAD ---------- */
type UserQuestionPayload struct {
	QuestionID int    `json:"question_id"`
	Title      string `json:"title"`
	UserQuery  string `json:"user_query"`
	// Optional: UserID if needed
}

/* ---------- RESPONSE ---------- */
type QuestionCheckResponse struct {
	Message       string      `json:"message"`
	Status        bool        `json:"status"`
	UserResult    interface{} `json:"user_result,omitempty"`
	CorrectResult interface{} `json:"correct_result,omitempty"`
}

func CheckUserAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("check user answer called")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	// Decode user payload
	fmt.Println(r.Body)
	var payload UserQuestionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Invalid request body",
			Status:  false,
		})
		return
	}

	fmt.Println("front pay", payload)

	/*
		User Payload Have
		 1. question title and id and user answer query ----> from Tales_Info table in Main Database
		 2. by question id and title we can find the question and used tables for tthis question
		 3. also get the correct annswer from main database

		 ----- Here all need of question done -----
		===========================================
		-------We Have to Run the user Query and Correct Answer Query-------

		after runing the both query we have to compare the results
		if both results are same then we have to return success response
		if both results are different then we have to return failure response with correct answer
		-----------------------------------------------------------
	*/

	// 1. Connect to main database to fetch question info
	clientMain := db.NewClient()
	if err := clientMain.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Main DB connection failed: " + err.Error(),
			Status:  false,
		})
		return
	}
	defer clientMain.Prisma.Disconnect()

	ctx := context.Background()

	// Fetch question by ID
	questionRecord, err := clientMain.QuestionRecords.FindUnique(
		db.QuestionRecords.ID.Equals(payload.QuestionID),
	).Exec(ctx)

	if err != nil || questionRecord == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Question not found",
			Status:  false,
		})
		return
	}

	// Get the correct answer query from DB
	answerVal, ok := questionRecord.Answer()
	var correctQuery string
	if !ok {
		correctQuery = ""
	} else {
		correctQuery = string(answerVal)
	}
	if correctQuery == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Correct answer not available for this question",
			Status:  false,
		})
		return
	}

	// 2. Connect to games database to execute queries
	clientGames := db2.NewClient()
	if err := clientGames.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Games DB connection failed: " + err.Error(),
			Status:  false,
		})
		return
	}
	defer clientGames.Prisma.Disconnect()

	// 3. Execute user's query
	var userResult []map[string]interface{}
	err = clientGames.Prisma.QueryRaw(payload.UserQuery).Exec(ctx, &userResult)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Error executing user query: " + err.Error(),
			Status:  false,
		})
		return
	}

	// 4. Execute correct answer query
	var correctResult []map[string]interface{}
	err = clientGames.Prisma.QueryRaw(correctQuery).Exec(ctx, &correctResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message: "Error executing correct answer query: " + err.Error(),
			Status:  false,
		})
		return
	}

	// 5. Compare results
	if reflect.DeepEqual(userResult, correctResult) {
		//Here u can make the question Status to completed for the user if needed
		//also update the user score if needed

		// Results match
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message:       "Correct! Your answer matches the expected result.",
			Status:        true,
			UserResult:    userResult,
			CorrectResult: correctResult,
		})
	} else {
		json.NewEncoder(w).Encode(QuestionCheckResponse{
			Message:       "Incorrect. Your answer does not match the expected result.",
			Status:        false,
			UserResult:    userResult,
			CorrectResult: correctResult,
		})
	}
}
