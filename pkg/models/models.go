package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no current command")

type Banner struct {
	ID           int                    `json:"id"`
	Tags         []int                  `json:"tag_ids"`
	FeatureId    int                    `json:"feature_id"`
	IsActive     bool                   `json:"is_active"`
	Content      map[string]interface{} `json:"content"`
	CreatedAt    time.Time              `json:"created_at"`
	LastModified time.Time              `json:"created_at"`
}
