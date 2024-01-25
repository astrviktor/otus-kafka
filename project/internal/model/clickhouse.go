package model

import "time"

type ClickhouseMessage struct {
	Id         string    `json:"id"`
	Status     JobStatus `json:"status"`
	CreateDate time.Time `json:"create_date"`
}
