package core

import (
	"github.com/gofrs/uuid"
)

func GenUID() string {
	uid, _ := uuid.NewV4()
	return uid.String()
}
