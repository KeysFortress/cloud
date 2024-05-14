package repositories

import (
	"time"

	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type PasswordRepository struct {
	Storage interfaces.Storage
}

func (pr *PasswordRepository) All() ([]models.Password, error) {
	query := `
		SELECT id,name, website, LENGTH(content) as password_lenght, created_at, updated_at
		FROM public.account_passwords
		GROUP BY id
	`

	queryResult := pr.Storage.Where(query, []interface{}{})

	var passwords []models.Password
	for queryResult.Next() {
		var password models.Password
		err := queryResult.Scan(
			&password.Id,
			&password.Email,
			&password.Website,
			&password.Password,
			&password.CreatedAt,
			&password.UpdatedAt,
		)

		if err != nil {
			return []models.Password{}, err
		}

		passwords = append(passwords, password)
	}

	return passwords, nil
}

func (pr *PasswordRepository) Add(passwordRequest dtos.IncomingPasswordRequest, id uuid.UUID) (uuid.UUID, error) {
	query := `
		INSERT INTO public.account_passwords(
		name, content, account_id, created_at,  website)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	queryResult := pr.Storage.Add(&query, &[]interface{}{
		&passwordRequest.Email,
		&passwordRequest.Password,
		&id,
		time.Now().UTC(),
		&passwordRequest.Website,
	})
	var createdId uuid.UUID
	err := queryResult.Scan(&createdId)

	if err != nil {
		return uuid.UUID{}, err
	}

	return createdId, nil
}

func (pr *PasswordRepository) Update(passwordRequest *dtos.IncomingPasswordUpdateRequest) bool {
	query := `
		UPDATE public.account_passwords
		SET name=$1, content=$2,  updated_at=$3, website=$4
		WHERE id=$5
	`

	result := pr.Storage.Exec(query, []interface{}{
		&passwordRequest.Email,
		&passwordRequest.Password,
		time.Now().UTC(),
		&passwordRequest.Website,
		&passwordRequest.Id,
	})

	return result
}

func (pr *PasswordRepository) Content(id uuid.UUID) (string, error) {
	query := `
		SELECT content
		FROM public.account_passwords
		WHERE id = $1
	`

	result := pr.Storage.Single(query, []interface{}{
		&id,
	})

	var passwordContent string
	err := result.Scan(&passwordContent)

	if err != nil {
		return "", err
	}

	return passwordContent, nil
}
