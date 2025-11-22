package routes

import (
	// "fmt"
	"gamecraft-backend/controllers/auths"
	"gamecraft-backend/controllers/question"
	"gamecraft-backend/middlewares"
	"net/http"
	"gamecraft-backend/controllers/sql"
)

func RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("/auth/signup", auths.SignUp)
	mux.HandleFunc("/auth/verify-otp", auths.VerifyOtp)
	mux.HandleFunc("/auth/login", auths.Login)
	mux.HandleFunc("/save-question", question.SaveQuestion)
	mux.HandleFunc("/add-game", sql.AddGame)
	mux.HandleFunc("/run",middlewares.AuthMiddleware(question.RunQuestion))
}

// only get request are Alloweed to this Function
func RegisterRouterGet(mux *http.ServeMux) {
	//here you can map the get requests to their respective handler functions

	mux.HandleFunc("/get-user", middlewares.AuthMiddleware(auths.GetUser))
	mux.HandleFunc("/logout", auths.Logout)
	mux.HandleFunc("/get-all-questions", question.GetAllQustion)
	mux.HandleFunc("/get-question", question.GetQustion)

	// fmt.Println("GET /getuser route registered")
}

// similar you can do for the put and delete reuest
