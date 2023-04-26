package errors

type InvalidParam struct {
	Param string `json:"param"`
}

func (i InvalidParam) Error() string {
	return "Invalid Parameter: " + i.Param
}
