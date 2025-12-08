package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	db2 "gamecraft-backend/prisma_games/prisma_games_client"
)

/* ---------- RESPONSE STRUCT ---------- */

func GetAllTablesPreview(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(GetTablesPreviewResponse{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	//there  is no params u have to send all tables present in the database

	// Connect to Games Database
	clientGames := db2.NewClient()
	if err := clientGames.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetTablesPreviewResponse{
			Message: "Database connection failed: " + err.Error(),
			Status:  false,
		})
		return
	}
	defer clientGames.Prisma.Disconnect()

	// Fetch all table names from the database
	var tableNames []map[string]interface{}
	if err := clientGames.Prisma.QueryRaw(`SELECT table_name FROM information_schema.tables WHERE table_schema='public';`).Exec(context.Background(), &tableNames); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetTablesPreviewResponse{
			Message: "Failed to fetch table names: " + err.Error(),
			Status:  false,
		})
		return
	}

	tablesPreview := make(map[string]interface{})

	// Iterate over each table and fetch preview data
	for _, row := range tableNames {
		tableName, ok := row["table_name"].(string)
		if !ok {
			// skip rows that don't have a valid table_name
			continue
		}

		// Fetch preview (first 5 rows) of the table
		var previewData []map[string]interface{}
		// quote table name to be safer with identifiers
		if err := clientGames.Prisma.QueryRaw(fmt.Sprintf(`SELECT * FROM "%s" LIMIT 5;`, tableName)).Exec(context.Background(), &previewData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(GetTablesPreviewResponse{
				Message: "Failed to fetch preview for table " + tableName + ": " + err.Error(),
				Status:  false,
			})
			return
		}

		tablesPreview[tableName] = previewData
	}

	// Send response
	json.NewEncoder(w).Encode(GetTablesPreviewResponse{
		Message: "Tables preview fetched successfully",
		Status:  true,
		Data:    tablesPreview,
	})
}