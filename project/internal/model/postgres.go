package model

type PostgresMessage struct {
	Schema  SchemaType      `json:"schema"`
	Payload PostgresPayload `json:"payload"`
}

type PostgresPayload struct {
	Id         string `json:"id"`
	Status     string `json:"status"`
	CreateDate int64  `json:"create_date"`
	FinishDate int64  `json:"finish_date"`
}

type SchemaType struct {
	Type   string      `json:"type"`
	Fields []FieldType `json:"fields"`
}

type FieldType struct {
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Field    string `json:"field"`
}

func (m *PostgresMessage) Init() {
	m.Schema = SchemaType{
		Type: "struct",
		Fields: []FieldType{
			{
				Field:    "id",
				Type:     "string",
				Optional: false,
			},
			{
				Field:    "status",
				Type:     "string",
				Optional: false,
			},
			{
				Field:    "create_date",
				Type:     "int64",
				Optional: false,
			},
			{
				Field:    "finish_date",
				Type:     "int64",
				Optional: false,
			},
		},
	}
}
