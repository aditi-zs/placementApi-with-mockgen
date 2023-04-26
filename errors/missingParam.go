package errors

import (
	"fmt"
	"strings"
)

type MissingParam struct {
	Param []string `json:"param"`
}

func (m MissingParam) Error() string {
	return fmt.Sprintf("Missing Parameter: %v", strings.Join(m.Param, ","))
}
