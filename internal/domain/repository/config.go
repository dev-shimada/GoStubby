package repository

import (
	"github.com/dev-shimada/gostubby/internal/domain/model"
)

type ConfigRepository interface {
	Load(path string) ([]model.Endpoint, error)
}
