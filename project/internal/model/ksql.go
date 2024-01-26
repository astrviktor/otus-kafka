package model

type KsqlData struct {
	Ksql              string            `json:"ksql"`
	StreamsProperties StreamsProperties `json:"streamsProperties"`
}

type StreamsProperties struct {
}

type KsqlResponse struct {
	Header       KsqlHeader
	Row          KsqlRow
	FinalMessage KsqlFinalMessage
}

type KsqlHeader struct {
	QueryId string `json:"queryId"`
	Schema  string `json:"schema"`
}

type KsqlRow struct {
	Row Row `json:"row"`
}

type Row struct {
	Columns []string `json:"columns"`
}

type KsqlFinalMessage struct {
	FinalMessage string `json:"finalMessage"`
}

type InformerResponse struct {
	Id         string `json:"id"`
	Status     string `json:"status"`
	CreateDate string `json:"create_date,omitempty"`
	FinishDate string `json:"finish_date,omitempty"`
}
