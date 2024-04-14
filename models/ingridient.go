package models

import "github.com/google/uuid"

type Ingridient struct {
	Id      uuid.UUID
	Name    string
	Kcal    float64
	Protein float64
	Carbs   float64
	Fat     float64
}
