package models

import (
	"github.com/google/uuid"
)

type Photo struct {
	UUID         uuid.UUID  `json:"id"`
	ItemUUID     *uuid.UUID `json:"item_id,omitempty"`
	PhotoPath    *string    `json:"photo_path,omitempty"`
	DisplayOrder *int32     `json:"display_order,omitempty"`
}

type PhotoID struct {
	UUID uuid.UUID `json:"id"`
}
