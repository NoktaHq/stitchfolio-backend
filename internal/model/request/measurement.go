package requestModel

import (
	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
)

type Measurement struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	Values entitiy_types.JSON `json:"values,omitempty"`

	PersonId    *uint `json:"personId,omitempty"`
	DressTypeId *uint `json:"dressTypeId,omitempty"`
	TakenById   *uint `json:"takenById,omitempty"`
}

type BulkMeasurementItem struct {
	DressTypeId uint               `json:"dressTypeId"`
	Values      entitiy_types.JSON `json:"values"`
}

type BulkMeasurementRequest struct {
	PersonId     uint                  `json:"personId"`
	Measurements []BulkMeasurementItem `json:"measurements"`
}
