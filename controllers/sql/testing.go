package sql

import (
    "encoding/json"
    "net/http"
    "time"
    "context"
    db "gamecraft-backend/prisma_games/prisma_games_client" // adjust import path
)

func AddGame(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(Response{
            Message: "Method not allowed",
            Status:  false,
        })
        return
    }

    // Parse request body
    var payload GamePayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(Response{
            Message: "Invalid request body",
            Status:  false,
        })
        return
    }

    // Connect Prisma client
    client := db.NewClient()
    if err := client.Connect(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(Response{
            Message: "Database connection failed: " + err.Error(),
            Status:  false,
        })
        return
    }
    defer client.Disconnect()

    // Insert new game
    ctx := context.Background()
    _, err := client.Game.CreateOne(
        db.Game.Title.Set(payload.Title),
        db.Game.Genre.Set(payload.Genre),
        db.Game.ReleaseDate.Set(time.Now()),
        db.Game.Developer.Set(payload.Developer),
    ).Exec(ctx)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(Response{
            Message: "Game creation failed: " + err.Error(),
            Status:  false,
        })
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{
        Message: "Game added successfully",
        Status:  true,
    })
}
