package uniquegen

import "github.com/google/uuid"

type UniqueGenerator string

func (g UniqueGenerator) NewId() string {
	return uuid.NewString()
}
