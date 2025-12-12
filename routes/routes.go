package routes

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"gamecraft-backend/controllers/auths"
	"gamecraft-backend/controllers/question"
	"gamecraft-backend/controllers/sql"
	"gamecraft-backend/controllers/users"
	"gamecraft-backend/middlewares"
	"net/http"
)

func RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("/auth/signup", auths.SignUp)
	mux.HandleFunc("/auth/verify-otp", auths.VerifyOtp)
	mux.HandleFunc("/auth/login", auths.Login)
	mux.HandleFunc("/save-question", question.SaveQuestion)
	mux.HandleFunc("/check-question", question.CheckQustion)
	mux.HandleFunc("/add-game", sql.AddGame)
	mux.HandleFunc("/run", middlewares.AuthMiddleware(question.RunQuestion))

	mux.HandleFunc("/create-table", question.CreateQuestionTable)
	mux.HandleFunc("/contribute-question", question.ContributeQuestion)
	mux.HandleFunc("/run-query", middlewares.AuthMiddleware(question.CheckUserAnswer))
	mux.HandleFunc("/update-question-solved-status", question.UpdateQuestionSolvedStatus)
	mux.HandleFunc("/update-user-profile", users.UpdateUserProfile)
}

// only get request are Alloweed to this Function
func RegisterRouterGet(mux *http.ServeMux) {
	//here you can map the get requests to their respective handler functions

	mux.HandleFunc("/get-user", middlewares.AuthMiddleware(auths.GetUser))
	mux.HandleFunc("/logout", auths.Logout)
	mux.HandleFunc("/get-all-questions", question.GetAllQustion)
	mux.HandleFunc("/get-question", question.GetQustion)
	mux.HandleFunc("/get-question-all", question.GetAllQustion)
	mux.HandleFunc("/get-tables-preview", question.GetTables)
	mux.HandleFunc("/get-all-tables-preview", question.GetAllTablesPreview)
	mux.HandleFunc("/get-questions-by-filters", question.GetQuestionsByFilters)
	mux.HandleFunc("/get-total-solved", users.GetTotalSolved)
	mux.HandleFunc("/get-user-profile", users.GetUserProfile)
	mux.HandleFunc("/getuser-activity", users.GetUserActivity)

	mux.HandleFunc("/datasharing", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Query())
		json.NewEncoder(w).Encode(auths.Response{
			Message: "questioned fetched successfully",
			Status:  true,
			Data: `{
    "open -na "Google Chrome" --args --disable-web-security --user-data-dir="/tmp/chrome"
"
}`,
		})
	})

	// fmt.Println("GET /getuser route registered")
}

// similar you can do for the put and delete reuest
