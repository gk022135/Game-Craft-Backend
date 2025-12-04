package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db "gamecraft-backend/prisma/db"
	db2 "gamecraft-backend/prisma_games/prisma_games_client"
)

/* ---------- REQUEST BODY STRUCT ---------- */

type CreateTableQuery struct {
	Query       string `json:"query"` // SQL query
	Owner       string `json:"owner"` // Optional
	TableName   string `json:"title"` // corrected
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"` // string works fine
}

/*
-------------Approach For this Controller-----------------
1. Validate HTTP method
2. Parse JSON body to get SQL query
3. Create Requested Table In Games Database Using ExecRaw
4. If Table Creation Successful, then maintaine the tables record in Main Database(Tables_Info)
5. If failed to creat tables then we return with server side error

-----------------------------------------------------------
*/

func CreateQuestionTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("creat table route url hitted")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	var requestData CreateTableQuery

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Invalid request body",
			Status:  false,
		})
		return
	}

	if requestData.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Message: "Query cannot be empty",
			Status:  false,
		})
		return
	}

	tableName := requestData.TableName
	createdBy := requestData.Owner
	description := requestData.Description
	query := requestData.Query

	/* ----------------- CONNECT TO GAMES DB ----------------- */

	gameClient := db2.NewClient()
	if err := gameClient.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Games Database connection failed",
			Status:  false,
		})
		return
	}
	defer gameClient.Prisma.Disconnect()

	// Correct raw SQL execution in prisma-client-go:
	_, err := gameClient.Prisma.ExecuteRaw(query).Exec(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: fmt.Sprintf("SQL Execution Error: %v", err),
			Status:  false,
		})
		return
	}

	/* ----------------- INSERT INTO MAIN DB ----------------- */

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Main Database connection failed",
			Status:  false,
		})
		return
	}
	defer client.Prisma.Disconnect()

	ctx := context.Background()

	_, err = client.TablesInfo.CreateOne(
		db.TablesInfo.TableName.Set(tableName),
		db.TablesInfo.Description.Set(description),
		db.TablesInfo.Querry.Set(query),
		db.TablesInfo.CreatedBy.Set(createdBy),
	).Exec(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Table record creation failed: " + err.Error(),
			Status:  false,
		})
		return
	}

	allTables, _ := client.TablesInfo.FindMany().Exec(ctx)

	/* ----------------- SUCCESS ----------------- */

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Table created successfully",
		"status":      true,
		"tables_info": allTables,
	})
}

// Code By Gkrrr
