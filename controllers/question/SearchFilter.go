package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"strconv"

	db "gamecraft-backend/prisma/db"
)

/* ---------- RESPONSE STRUCT ---------- */

type FilteredQuestionResponse struct {
	Id              int      `json:"Id"`
	Title           string   `json:"Title"`
	Description     string   `json:"Description"`
	Topics          []string `json:"Topics"`
	UsedTables      []string `json:"UsedTables"`
	ContributedBy   string   `json:"ContributedBy"`
	Points          int      `json:"Points,omitempty"`
	Answer          string   `json:"Answer,omitempty"`
	DifficultyLevel string   `json:"DifficultyLevel"`
}

/* ---------- HANDLER ---------- */

func GetQuestionsByFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	// Query params
	search := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("search")))
	level := strings.TrimSpace(r.URL.Query().Get("level"))

	// Prisma Client
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "failed to connect to server", http.StatusInternalServerError)
		return
	}
	defer client.Prisma.Disconnect()

	ctx := context.Background()

	/* ----------------- DYNAMIC FILTERS ----------------- */

	filters := []db.QuestionRecordsWhereParam{}

	// Difficulty level filter
	if level != "" {
		filters = append(filters,
			db.QuestionRecords.DifficultyLevel.Equals(level),
		)
	}

	// Search filter across multiple fields
	// If search is numeric → convert to int
	var searchID int
	isNumeric := false

	if idVal, err := strconv.Atoi(search); err == nil {
		searchID = idVal
		isNumeric = true
	}

	if search != "" {
		orFilters := []db.QuestionRecordsWhereParam{
			db.QuestionRecords.Title.Contains(search),
			db.QuestionRecords.Description.Contains(search),
			db.QuestionRecords.ContributedBy.Contains(search),
			db.QuestionRecords.Topics.Has(search),
			db.QuestionRecords.UsedTables.Has(search),
		}

		// If user searched for a number → search in ID
		if isNumeric {
			orFilters = append(orFilters,
				db.QuestionRecords.ID.Equals(searchID),
			)
		}

		filters = append(filters, db.QuestionRecords.Or(orFilters...))
	}

	/* ----------------- QUERY DATABASE ----------------- */

	// fmt.Println("Fetching level Medium questions...")
	// fmt.Println("search:", search, " level:", level)

	// results, err := client.QuestionRecords.FindMany(
	// 	db.QuestionRecords.Title.Contains("Get all orders with user details"),
	// ).Exec(context.Background())

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("Found:", results)
	// }

	questions, err := client.QuestionRecords.
		FindMany(filters...).
		OrderBy(db.QuestionRecords.ID.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Failed to fetch filtered questions: " + err.Error(),
			Status:  false,
		})
		return
	}

	fmt.Println("Filtered questions fetched:", len(questions))

	/* ----------------- MAP RESULTS ----------------- */

	var result []FilteredQuestionResponse

	for _, q := range questions {

		points := 0
		if val, ok := q.Rewards(); ok {
			points = int(val)
		}

		answer := ""
		if val, ok := q.Answer(); ok {
			answer = val
		}

		result = append(result, FilteredQuestionResponse{
			Id:              q.ID,
			Title:           q.Title,
			Description:     q.Description,
			Topics:          q.Topics,
			UsedTables:      q.UsedTables,
			ContributedBy:   q.ContributedBy,
			Points:          points,
			Answer:          answer,
			DifficultyLevel: q.DifficultyLevel,
		})
	}

	/* ----------------- SEND RESPONSE ----------------- */

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Filtered questions fetched successfully",
		Status:  true,
		Data:    result,
	})
}
