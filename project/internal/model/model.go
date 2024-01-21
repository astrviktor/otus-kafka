package model

import "time"

type Job struct {
	Id         string    `json:"id"`
	DateCreate time.Time `json:"date_create"`
	Sleep      int64     `json:"sleep"`
	Data       string    `json:"data"`
}

type JobResponse struct {
	Id         string    `json:"id"`
	DateCreate time.Time `json:"date_create"`
	Sleep      int64     `json:"sleep"`
	DataSize   int       `json:"data_size"`
}

type JobDone struct {
	Id         string    `json:"id"`
	Status     string    `json:"status"`
	DateCreate time.Time `json:"date_create"`
	DateDone   time.Time `json:"date_done"`
}

type PostgresMessage struct {
	Id         string    `json:"id"`
	Status     string    `json:"status"`
	DateCreate time.Time `json:"date_create"`
	DateDone   time.Time `json:"date_done"`
}

type ClickhouseMessage struct {
	Id         string    `json:"id"`
	Status     string    `json:"status"`
	DateCreate time.Time `json:"date_create"`
	DateDone   time.Time `json:"date_done"`
}
