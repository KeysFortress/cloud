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

func (pr *SecretsRepository) Get(id uuid.UUID) (models.Secret, error) {
	query := `
		SELECT id,name, website, LENGTH(content) as password_lenght, created_at, updated_at
		FROM public.account_secrets
		WHERE id=$1
	`

	queryResult := pr.Storage.Single(query, []interface{}{
		&id,
	})
	var password models.Secret
	err := queryResult.Scan(
		&password.Id,
		&password.Email,
		&password.Website,
		&password.Password,
		&password.CreatedAt,
		&password.UpdatedAt,
	)

	if err != nil {
		return models.Secret{}, err
	}

	return password, nil
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
		&incomingSecret.Secret,
		&id,
		time.Now().UTC(),
		&incomingSecret.Website,
	})
	var createdId uuid.UUID
	err := queryResult.Scan(&createdId)

	if err != nil {
		return uuid.UUID{}, err
	}

	return createdId, nil

}

func (sr *SecretsRepository) Update(secretRequest *dtos.IncomingSecretsUpdateRequest) bool {
	query := `
		UPDATE public.account_secrets
		SET name=$1, content=$2,  updated_at=$3, website=$4, description = $5
		WHERE id=$6;
	`

	result := sr.Storage.Exec(query, []interface{}{
		&secretRequest.Email,
		&secretRequest.Secret,
		time.Now().UTC(),
		&secretRequest.Website,
		"--",
		&secretRequest.Id,
	})

	return result

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

func (sr *SecretsRepository) Delete(id uuid.UUID) bool {

	query := `
		DELETE FROM public.account_secrets
		WHERE id = $1
	`
	result := sr.Storage.Exec(query, []interface{}{
		&id,
	})

	return result

}
