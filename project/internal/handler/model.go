package handler

type Job struct {
	Id    string `json:"id"`
	Sleep int64  `json:"sleep"`
	Data  string `json:"data"`
}

type JobResponse struct {
	Id       string `json:"id"`
	Sleep    int64  `json:"sleep"`
	DataSize int    `json:"data_size"`
}
