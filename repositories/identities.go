package repositories

import (
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type IdentityRepository struct {
	Storage interfaces.Storage
}

func (i *IdentityRepository) All() ([]models.Identity, error) {
	return []models.Identity{}, nil
}

func (i *IdentityRepository) Get(id uuid.UUID) (models.Identity, error) {
	return models.Identity{}, nil
}

func (i *IdentityRepository) GetKeyTypes() ([]models.KeyType, error) {
	return []models.KeyType{}, nil
}

func (i *IdentityRepository) Add(identity dtos.CreateIdentity, keyData any) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (i *IdentityRepository) Update(identity dtos.UpdateIdentity, keyData any) bool {
	return true
}
