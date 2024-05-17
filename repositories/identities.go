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

func (i *IdentityRepository) GetInternal(id uuid.UUID) (models.IdentityInternal, error) {
	return models.IdentityInternal{}, nil
}

func (i *IdentityRepository) GetKeyTypes() ([]models.KeyType, error) {
	return []models.KeyType{}, nil
}

func (i *IdentityRepository) Add(identity dtos.CreateIdentity, public *[]byte, private *[]byte) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (i *IdentityRepository) Update(identity *dtos.UpdateIdentity, public *[]byte, private *[]byte) bool {
	return true
}
