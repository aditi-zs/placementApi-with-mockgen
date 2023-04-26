package entities

import "github.com/google/uuid"

type Company struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Category Category  `json:"category,omitempty"`
}

type Category string

const (
	MASS      Category = "MASS"
	DREAMIT   Category = "DREAM IT"
	OPENDREAM Category = "OPEN DREAM"
	CORE      Category = "CORE"
)

func IsValidCategory(c Category) bool {
	if c == MASS || c == DREAMIT || c == OPENDREAM || c == CORE {
		return true
	}

	return false
}
