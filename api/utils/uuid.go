package utils

import (
	uuid "github.com/satori/go.uuid"
)

func UUIDNew() string {
	u1 := uuid.NewV4()
	return u1.String()
}
