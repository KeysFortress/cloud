package repositories

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type IdentityRepository struct {
	Storage interfaces.Storage
}

func (i *IdentityRepository) All() ([]models.Identity, error) {
	query := `
		SELECT ai.name,  kt.name, ai.key_size, ai.identity_key,LENGTH(ai.identity_secret_key), ai.created_at, ai.updated_at FROM public.account_identities as ai
		JOIN public.key_types as  kt on kt.id = ai.key_type_id
	`

	queryResult := i.Storage.Where(query, []interface{}{})

	var identities []models.Identity
	for queryResult.Next() {
		var identity models.Identity

		err := queryResult.Scan(
			&identity.Id,
			&identity.KeyType,
			&identity.KeySize,
			&identity.PublicKey,
			&identity.PrivateKey,
			&identity.CreatedAt,
			&identity.UpdatedAt,
		)

		if err != nil {
			fmt.Println(err)
			return []models.Identity{}, err
		}

		identities = append(identities, identity)
	}

	return identities, nil
}

func (i *IdentityRepository) Get(id uuid.UUID) (models.Identity, error) {
	query := `
		SELECT ai.name,  kt.name, ai.key_size, ai.identity_key,LENGTH(ai.identity_secret_key), ai.created_at, ai.updated_at FROM public.account_identities as ai
		JOIN public.key_types as  kt on kt.id = ai.key_type_id
		WHERE ai.id = $1
	`

	queryResult := i.Storage.Single(query, []interface{}{
		&id,
	})

	var identity models.Identity
	err := queryResult.Scan(
		&identity.Id,
		&identity.KeyType,
		&identity.KeySize,
		&identity.PublicKey,
		&identity.PrivateKey,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err)
		return models.Identity{}, err
	}

	return identity, nil
}

func (i *IdentityRepository) GetInternal(id uuid.UUID) (models.IdentityInternal, error) {
	query := `
		SELECT ai.name, ai.key_type_id, ai.key_size, ai.identity_key,ai.identity_secret_key, ai.created_at, ai.updated_at FROM public.account_identities as ai
		WHERE ai.id = $1
	`

	queryResult := i.Storage.Single(query, []interface{}{
		&id,
	})

	var identity models.IdentityInternal
	err := queryResult.Scan(
		&identity.Id,
		&identity.KeyType,
		&identity.KeySize,
		&identity.PublicKey,
		&identity.PrivateKey,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err)
		return models.IdentityInternal{}, err
	}

	return identity, nil
}

func (i *IdentityRepository) GetKeyTypes() ([]models.KeyType, error) {
	query := `
		SELECT * FROM public.key_types
	`

	queryResult := i.Storage.Where(query, []interface{}{})

	var types []models.KeyType
	for queryResult.Next() {
		var keyType models.KeyType

		err := queryResult.Scan(
			&keyType.Id,
			&keyType.Name,
			&keyType.Description,
			&keyType.HasSize,
		)

		if err != nil {
			fmt.Println(err)
			return []models.KeyType{}, err
		}

		types = append(types, keyType)
	}

	return types, nil
}

func (i *IdentityRepository) Add(identity dtos.CreateIdentity, public *[]byte, private *[]byte, id uuid.UUID) (uuid.UUID, error) {
	query := `
		INSERT INTO public.account_identities(
		name, identity_key, identity_secret_key, key_size, key_type_id, account_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	queryResult := i.Storage.Single(query, []interface{}{
		identity.Name,
		base64.StdEncoding.EncodeToString(*public),
		base64.StdEncoding.EncodeToString(*private),
		identity.KeyType,
		identity.KeySize,
		id,
		time.Now().UTC(),
	})

	var created uuid.UUID
	err := queryResult.Scan(&created)

	if err != nil {
		fmt.Println(err)
		return uuid.UUID{}, err
	}

	return created, nil
}

func (i *IdentityRepository) Update(identity *dtos.UpdateIdentity, public *[]byte, private *[]byte) bool {
	query := `
		UPDATE public.account_identities
		SET
		 	name=$1,
			identity_key=$2,
			identity_secret_key=$3,
			key_size=$4,
			key_type_id=$5
			updated_at=$6,
		WHERE id=$7;
	`

	updated := i.Storage.Exec(query, []interface{}{
		identity.Name,
		base64.StdEncoding.EncodeToString(*public),
		base64.StdEncoding.EncodeToString(*private),
		identity.KeySize,
		identity.KeyType,
		time.Now().UTC(),
		identity.Id,
	})

	return updated
}
