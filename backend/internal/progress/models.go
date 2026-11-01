package progress

import "github.com/THARUN-BART/LearnX/backend/internal/dsa/model"

type Record = model.Progress

type UpdateRequest struct {
	Status string `json:"status"`
	Reset  bool   `json:"reset"`
}
