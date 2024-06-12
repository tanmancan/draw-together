package helpers

import (
	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/model"
)

func ProtoToUUID(id *model.UUID) (uuid.UUID, error) {
	return uuid.Parse(id.Value)
}

func ProtoFromUUID(uid uuid.UUID) *model.UUID {
	id := &model.UUID{
		Value: uid.String(),
	}
	return id
}
