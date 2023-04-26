package errors

type DB struct {
	Reason string `json:"reason"`
}

func (d DB) Error() string {
	return "DB Error: " + d.Reason
}
