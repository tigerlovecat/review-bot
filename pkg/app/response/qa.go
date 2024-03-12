package response

type QaListTestResponse struct {
	Time     string `json:"time"`
	UserName string `json:"username"`
	Message  string `json:"message"`
	Reply    string `json:"reply"`
}
