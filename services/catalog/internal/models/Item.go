package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	Category    string     `json:"category"`
	Price       int32      `json:"price"`
	Status      string     `json:"status"`
	Brand       string     `json:"brand"`
	Color       string     `json:"color"`
	Size        string     `json:"size"`
	Sex         string     `json:"sex"`
	Description *string    `json:"description,omitempty"`
	Created_at  *time.Time `json:"created_at,omitempty"`
	Sold_at     *time.Time `json:"sold_at,omitempty"`
	Seller_id   uuid.UUID  `json:"seller_id"`
	Buyer_id    *uuid.UUID `json:"buyer_id,omitempty"`
}

type ItemUpdate struct {
	UUID        uuid.UUID  `json:"uuid"`
	Category    *string    `json:"category,omitempty"`
	Price       *int32     `json:"price,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Brand       *string    `json:"brand,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Size        *string    `json:"size,omitempty"`
	Sex         *string    `json:"sex,omitempty"`
	Description *string    `json:"description,omitempty"`
	Created_at  *time.Time `json:"created_at,omitempty"`
	Sold_at     *time.Time `json:"sold_at,omitempty"`
	Seller_id   *uuid.UUID `json:"seller_id,omitempty"`
	Buyer_id    *uuid.UUID `json:"buyer_id,omitempty"`
}

type ItemID struct {
	UUID uuid.UUID `json:"uuid"`
}
