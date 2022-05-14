package dto

import "time"

type Name struct {
	First  string `json:"first"`
	Middle string `json:"middle"`
	Last   string `json:"last"`
}

type SourceDto struct {
	Name      Name                   `json:"name"`
	Age       uint32                 `json:"age"`
	Friends   []Name                 `json:"friends"`
	CreatedAt time.Time              `json:"created_at"`
	Data      map[string]interface{} `json:"data"`
}
