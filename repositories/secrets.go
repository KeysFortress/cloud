package repositories

import (
	"time"

	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type SecretsRepository struct {
	Storage interfaces.Storage
}

func (pr *SecretsRepository) All() ([]models.Secret, error) {
	query := `
		SELECT id,name, description, website, LENGTH(content) as password_lenght, created_at, updated_at
		FROM public.account_secrets
		GROUP BY id
	`

	queryResult := pr.Storage.Where(query, []interface{}{})

	var secrets []models.Secret
	for queryResult.Next() {
		var secret models.Secret
		err := queryResult.Scan(
			&secret.Id,
			&secret.Email,
			&secret.Description,
			&secret.Website,
			&secret.Password,
			&secret.CreatedAt,
			&secret.UpdatedAt,
		)

		if err != nil {
			return []models.Secret{}, err
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (sr *SecretsRepository) Add(incomingSecret dtos.IncomingSecretsRequest, id uuid.UUID) (uuid.UUID, error) {
	query := `
		INSERT INTO public.account_secrets(
		name, description, content, account_id, created_at,  website)
		VALUES ($1, $2,$3, $4, $5, $6)
		RETURNING id
	`

	queryResult := sr.Storage.Add(&query, &[]interface{}{
		&incomingSecret.Email,
		&incomingSecret.Description,
		&incomingSecret.Password,
		&id,
		&time.UTC,
		&incomingSecret.Website,
	})
	var createdId uuid.UUID
	err := queryResult.Scan(&createdId)

	if err != nil {
		return uuid.UUID{}, err
	}

	return createdId, nil

}

func (sr *SecretsRepository) Update(secretRequest models.Secret) (bool, error) {
	query := `
		UPDATE public.account_secrets
		SET name=$1, content=$2,  updated_at=$3, website=$4, description = $6
		WHERE id=$7;
	`

	result := sr.Storage.Single(query, []interface{}{
		&secretRequest.Email,
		&secretRequest.Password,
		&time.UTC,
		&secretRequest.Website,
		&secretRequest.Description,
		&secretRequest.Id,
	})

	err := result.Scan()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (sr *SecretsRepository) Content(id uuid.UUID) (string, error) {
	query := `
		SELECT content
		FROM public.account_secrets
		WHERE id = $1
	`

	result := sr.Storage.Single(query, []interface{}{
		&id,
	})

	var passwordContent string
	err := result.Scan(&passwordContent)

	if err != nil {
		return "", err
	}

	return passwordContent, nil
}
