package models

import "errors"

var ErrNoRecord = errors.New("models: no current command")

type Banner struct {
	ID         int                    `json:"id"`
	Tags       []int                  `json:"tag_ids"`
	Feature_id int                    `json:"feature_id"`
	Is_active  bool                   `json:"is_active"`
	Content    map[string]interface{} `json:"content"`
}
