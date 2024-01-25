package model

import "time"

type Job struct {
	Id         string    `json:"id"`
	CreateDate time.Time `json:"create_date"`
	Sleep      int64     `json:"sleep"`
	Data       string    `json:"data"`
}

type JobResponse struct {
	Id         string    `json:"id"`
	CreateDate time.Time `json:"create_date"`
	Sleep      int64     `json:"sleep"`
	DataSize   int       `json:"data_size"`
}

type JobStatus string

const (
	JobStatusCreate JobStatus = "create"
	JobStatusFinish JobStatus = "finish"
)

type JobDone struct {
	Id         string    `json:"id"`
	Status     JobStatus `json:"status"`
	CreateDate time.Time `json:"create_date"`
	FinishDate time.Time `json:"finish_date"`
}
