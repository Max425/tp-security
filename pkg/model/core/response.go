package core

type Response struct {
	ID      string            `json:"id" bson:"_id"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
