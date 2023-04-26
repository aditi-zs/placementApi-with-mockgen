package entities

import "github.com/google/uuid"

type Student struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Phone  string    `json:"phone"`
	DOB    string    `json:"dob"`
	Branch Branch    `json:"branch"`
	Comp   Company   `json:"comp,omitempty"`
	Status Status    `json:"status"`
}

type Branch string
type Status string

const (
	CSE   Branch = "CSE"
	ISE   Branch = "ISE"
	MECH  Branch = "MECH"
	ECE   Branch = "ECE"
	EEE   Branch = "EEE"
	CIVIL Branch = "CIVIL"
)

const (
	ACCEPTED Status = "ACCEPTED"
	REJECTED Status = "REJECTED"
	PENDING  Status = "PENDING"
)

func IsValidStatus(s Status) bool {
	if s == ACCEPTED || s == REJECTED || s == PENDING {
		return true
	}

	return false
}

func IsValidBranch(b Branch) bool {
	switch b {
	case CSE, ISE, MECH, ECE, EEE, CIVIL:
		return true
	default:
		return false
	}
}
