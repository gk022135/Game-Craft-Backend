package sql

type Response struct {
    Message string `json:"message"`
    Status  bool   `json:"status"`
}

type GamePayload struct {
    Title     string `json:"title"`
    Genre     string `json:"genre"`
    Developer string `json:"developer"`
}