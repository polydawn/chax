package protocol

import (
	"fmt"
)

type ServerDesc struct {
	Host string
	Port uint16
}

func (s *ServerDesc) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
