package models

type Pool struct {
	ID        int64    `json:"id"`
	Name      int64    `json:"name"`
	Input     string   `json:"input"`
	Output    string   `json:"output"`
	ProjectId int64    `json:"userId"`
	Project   *Project `pg:"rel:has-one" json:"project"`
}
