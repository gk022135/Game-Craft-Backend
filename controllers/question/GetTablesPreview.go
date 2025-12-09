package question

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db2 "gamecraft-backend/prisma_games/prisma_games_client"
)

/* ---------- RESPONSE STRUCT ---------- */

type GetTablesPreviewResponse struct {
	Message string                 `json:"message"`
	Status  bool                   `json:"status"`
	Data    map[string]interface{} `json:"data"` // tableName → rows preview
}

func GetTables(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(GetTablesPreviewResponse{
			Message: "Method not allowed",
			Status:  false,
		})
		return
	}

	// Get multiple tables from query param: ?tables=users,orders,products
	rawTables := r.URL.Query().Get("tables")
	if rawTables == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GetTablesPreviewResponse{
			Message: "No table names provided",
			Status:  false,
		})
		return
	}

	// Convert comma-separated string → array
	tables := strings.Split(rawTables, ",")

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


	// Fetch preview data for each table----> make the required data
	ctx := context.Background()
	tablesData := make(map[string]interface{})


//------------ Important part ---------//
	// Fetch 3–5 rows from each table
	for _, table := range tables {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}
		// fmt.Print("hi ", table)

		query := fmt.Sprintf("SELECT * FROM %s LIMIT 5;", table)

		var result []map[string]interface{}

		// ExecRawQuery is required to bind results from raw SQL
		err := clientGames.Prisma.QueryRaw(query).Exec(ctx, &result)
		if err != nil {
			tablesData[table] = map[string]string{
				"error": "Failed to fetch preview: " + err.Error(),
			}
			continue
		}

		tablesData[table] = result
	}
//------------ Important part ends ---------//


	// Response
	if err := json.NewEncoder(w).Encode(GetTablesPreviewResponse{
		Message: "Tables preview fetched successfully",
		Status:  true,
		Data:    tablesData,
	}); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
