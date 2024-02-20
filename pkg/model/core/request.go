package core

type Request struct {
	ID         string            `json:"id" bson:"_id"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	GetParams  map[string]string `json:"get_params"`
	Headers    map[string]string `json:"headers"`
	Cookies    map[string]string `json:"cookies"`
	PostParams map[string]string `json:"post_params"`
	Response   Response          `json:"response"`
}
