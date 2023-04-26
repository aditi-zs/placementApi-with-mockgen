package errors

type EntityNotFound struct {
	Reason string
}

func (e EntityNotFound) Error() string {
	return "Entity Not Found:" + e.Reason
}
