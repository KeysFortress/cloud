package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type IngridientRepository struct {
	Storage interfaces.Storage
}

func (ingridientRepository *IngridientRepository) OpenConnection(storage *interfaces.Storage) bool {

	ingridientRepository.Storage = *storage
	ingridientRepository.Storage.Open()
	return true
}

func (ingridientRepository *IngridientRepository) Add(ingridient *models.Ingridient) (*uuid.UUID, error) {

	var query = `
		INSERT INTO public.ingridients(
		name, description, created_at, fat, carbs, protein, kcal)
		VALUES ($1,'--', $2, $3, $4, $5, $6)
		RETURNING id

	`
	queryResult := ingridientRepository.Storage.Add(&query, &[]interface{}{
		&ingridient.Name,
		time.Now().UTC(),
		&ingridient.Fat,
		&ingridient.Carbs,
		&ingridient.Protein,
		&ingridient.Kcal,
	})
	var id uuid.UUID
	err := queryResult.Scan(&id)

	if err != nil {
		fmt.Println("Failed to create a new entry for ingridient")
		fmt.Println(ingridient)
		return nil, err
	}

	return &id, nil
}

func (ingridientRepository *IngridientRepository) Close() {
	ingridientRepository.Storage.Close()
}
