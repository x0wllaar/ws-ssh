package util

import "github.com/google/uuid"

func StringGuid() string {
	id := uuid.New()
	return id.String()
}
