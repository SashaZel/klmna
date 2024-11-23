package api

import (
	"klmna/internal/db/models"
)

type TaskResponse struct {
	Ok       bool              `json:"ok"`
	Error    string            `json:"error"`
	Input    *models.TaskInput `json:"input"`
	Template string            `json:"template"`
}
