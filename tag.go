package seekpo

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID
	Name string
	Unit string
	Type Type
}
