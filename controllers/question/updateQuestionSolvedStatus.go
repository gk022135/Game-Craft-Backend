package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "gamecraft-backend/prisma/db"
)

func UpdateQuestionSolvedStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	// <---------- Read query params --------->
	userIDStr := r.URL.Query().Get("user_id") //---->> should be email
	questionIDStr := r.URL.Query().Get("question_id")
	fmt.Println("useremail", userIDStr, questionIDStr)

	// Convert params
	// userID, err := strconv.Atoi(userIDStr)
	// if err != nil {
	// 	http.Error(w, "invalid user_id", http.StatusBadRequest)
	// 	return
	// }

	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, "invalid question_id", http.StatusBadRequest)
		return
	}



	// <---------- DB connect --------->
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "failed to connect to DB", http.StatusInternalServerError)
		return
	}
	defer client.Prisma.Disconnect()

	//<--------- Get user --------->
	user, err := client.User.FindUnique(
		db.User.Email.Equals(userIDStr),
	).Exec(context.Background())

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Message: "User not found",
			Status:  false,
		})
		return
	}

	//<------------ Check if already solved  ----------->
	if containsInt(user.SolvedEasy, questionID) ||
		containsInt(user.SolvedMedium, questionID) ||
		containsInt(user.SolvedHard, questionID) {

		json.NewEncoder(w).Encode(Response{
			Message: "Question already solved",
			Status:  true,
			Data: map[string]interface{}{
				"earned_points": user.EarnedPoints,
			},
		})
		return
	}

	// <------------- Get question Details -------------->
	question, err := client.QuestionRecords.FindUnique(
		db.QuestionRecords.ID.Equals(questionID),
	).Exec(context.Background())

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Message: "Question not found",
			Status:  false,
		})
		return
	}


	//<---------- Extract points and difficulty ---------->
	points := question.Rewards
	qType := question.DifficultyLevel

	// assume EarnedPoints is an int field on the user struct
	p, ok := points()
	var pointsInt int

	if ok {
		pointsInt = int(p) // convert db.Int → int
	} else {
		pointsInt = 0 // default value if undefined
	}
	currentPoints := user.EarnedPoints
	newPoints := currentPoints + pointsInt

	// newPoints := currentPoints + points ---> why do not use points directly
	// because points is of type db.Int which is a pointer type(function) in prisma go client
	// <----------- End Of Points Extraction ----------->




	// <----------- Update user record (Update based on difficulty)----------->
	var updateErr error
	switch qType {
	case "Easy":
		_, updateErr = client.User.FindUnique(
			db.User.Email.Equals(userIDStr),
		).Update(
			db.User.SolvedEasy.Push([]int{questionID}),
			db.User.EarnedPoints.Set(newPoints),
		).Exec(context.Background())

	case "Medium":
		_, updateErr = client.User.FindUnique(
			db.User.Email.Equals(userIDStr),
		).Update(
			db.User.SolvedMedium.Push([]int{questionID}),
			db.User.EarnedPoints.Set(newPoints),
		).Exec(context.Background())

	case "Hard":
		_, updateErr = client.User.FindUnique(
			db.User.Email.Equals(userIDStr),
		).Update(
			db.User.SolvedHard.Push([]int{questionID}),
			db.User.EarnedPoints.Set(newPoints),
		).Exec(context.Background())

	default:
		json.NewEncoder(w).Encode(Response{
			Message: "Unknown difficulty level",
			Status:  false,
		})
		return
	}

	if updateErr != nil {
		json.NewEncoder(w).Encode(Response{
			Message: "Failed to update user",
			Status:  false,
		})
		return
	}


	//<----------- Update successful response ----------->

	// Success
	json.NewEncoder(w).Encode(Response{
		Message: "Question status updated successfully",
		Status:  true,
		Data: map[string]interface{}{
			"earned_points": newPoints,
		},
	})
}

func containsInt(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
