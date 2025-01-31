package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Sex         string `json:"sex"`
	Category    string `json:"category"`
	Brand       string `json:"brand"`
	Color       string `json:"color"`
	Size        string `json:"size"`
	Description string `json:"description,omitempty"`
	Price       int32  `json:"price"`
}

type ItemUpdate struct {
	UUID        uuid.UUID `json:"uuid"`
	Sex         *string   `json:"sex,omitempty"`
	Category    *string   `json:"category,omitempty"`
	Brand       *string   `json:"brand,omitempty"`
	Color       *string   `json:"color,omitempty"`
	Size        *string   `json:"size,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *int32    `json:"price,omitempty"`
}

type ItemID struct {
	UUID uuid.UUID `json:"uuid"`
}
