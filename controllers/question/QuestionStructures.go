package question

type QuestionController struct {
	Title string `json:"Title"`
	Description string `json:"Description"`
	StarterSchema  string `json:"StarterSchema"`
	StarterData     string `json:"StarterData"`
	CorrectQuery  string `json:"CorrectQuery"`
}

type Response struct {
	Message  string      `json:"message"`
	Status   bool        `json:"status"`
	TryLater string      `json:"try_later,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type QuestionResponse struct {
	Id	int	`json:"Id"`
	Title string `json:"Title"`
	Description string `json:"Description"`
}

type SolutionController struct {
	AnswerQuery  string `json:"AnswerQuery"`
}

type SolutionResponse struct {
	IsCorrect  string `json:"IsCorrect"`
	Message  string `json:"Message"`
}