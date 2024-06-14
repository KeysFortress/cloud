package repositories

import (
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
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

func (er *EventRepository) TakeCount(rangeValue *dtos.ValueRage) ([]models.Event, error) {
	sql := `
		SELECT * FROM events
		ORDER BY created_at
		LIMIT $1 OFFSET $2
	`

	query := er.Storage.Where(sql, []interface{}{
		&rangeValue.Take,
		&rangeValue.Skip,
	})

	var events []models.Event
	for query.Next() {
		var event models.Event
		err := query.Scan(&event.Id, &event.Description, &event.Device, &event.EventDate, &event.TypeId)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}

func (er *EventRepository) GetByType(eventType int) ([]models.Event, error) {
	sql := `
		SELECT * FROM events
		WHERE event_type_id = $1
	`

	query := er.Storage.Where(sql, []interface{}{
		&eventType,
	})

	var events []models.Event
	for query.Next() {
		var event models.Event
		err := query.Scan(&event.Id, &event.Description, &event.Device, &event.EventDate, &event.TypeId)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
