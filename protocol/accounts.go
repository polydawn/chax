package protocol

import (
	"fmt"
)

type Account struct {
	Username string
	Domain   string
	Password string
}

func (a Account) FQAN() string {
	return fmt.Sprintf("%s@%s", a.Username, a.Domain)
}
