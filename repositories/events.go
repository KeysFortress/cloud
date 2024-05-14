package repositories

import (
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
)

type EventRepository struct {
	Storage interfaces.Storage
}

func (er *EventRepository) Add(event dtos.CreateEvent) (uuid.UUID, error) {
	sql := `
	INSERT INTO public.events(
		event_type_id, description, device_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	result := er.Storage.Add(&sql, &[]interface{}{
		event.TypeId,
		event.Description,
		event.DeviceId,
		event.CreatedAt,
	})

	var id uuid.UUID
	err := result.Scan(&id)

	if err != nil {
		return uuid.UUID{}, nil
	}

	return id, nil
}
