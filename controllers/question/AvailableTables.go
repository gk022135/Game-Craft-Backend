package question

import (
	"context"
	"encoding/json"
	"net/http"

	db2 "gamecraft-backend/prisma_games/prisma_games_client"
)

/* ---------- RESPONSE STRUCT ---------- */

type GetAllTablesResponse struct {
	Message string        `json:"message"`
	Status  bool          `json:"status"`
	Tables  []string      `json:"tables"` // list of table names
}

func GetAllTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(GetAllTablesResponse{
			Message: "Method not allowed",
			Status:  false,
			Tables:  nil,
		})
		return
	}

	// Connect to Games Database
	clientGames := db2.NewClient()
	if err := clientGames.Prisma.Connect(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetAllTablesResponse{
			Message: "Database connection failed: " + err.Error(),
			Status:  false,
			Tables:  nil,
		})
		return
	}
	defer clientGames.Prisma.Disconnect()

	ctx := context.Background()

	// Query database schema to get all table names
	var tables []string
	err := clientGames.Prisma.QueryRaw("SHOW TABLES").Exec(ctx, &tables)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GetAllTablesResponse{
			Message: "Failed to fetch tables: " + err.Error(),
			Status:  false,
			Tables:  nil,
		})
		return
	}

	// Success response
	json.NewEncoder(w).Encode(GetAllTablesResponse{
		Message: "Tables fetched successfully",
		Status:  true,
		Tables:  tables,
	})
}
