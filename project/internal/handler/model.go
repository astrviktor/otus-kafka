package handler

type Job struct {
	Id    string `json:"id"`
	Sleep int64  `json:"sleep"`
	Data  string `json:"data"`
}
