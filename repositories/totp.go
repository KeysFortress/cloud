package repositories

import (
	"time"

	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type TotpRepository struct {
	Storage interfaces.Storage
}

func (tr *TotpRepository) All() ([]models.TimeBasedCode, error) {
	query := `
		SELECT t.id, t.website, t.email, LENGTH(t.secret) as password_lenght, tt.name, t.expiration, t.created_at, t.updated_at
		FROM public.time_based_codes as t
		JOIN public.time_based_code_types as tt
		on tt.id = t.type_id
 	`

	queryResult := tr.Storage.Where(query, []interface{}{})
	if queryResult == nil {
		return []models.TimeBasedCode{}, nil
	}
	var timeBasedCodes []models.TimeBasedCode

	for queryResult.Next() {
		var timeBasedCode models.TimeBasedCode
		err := queryResult.Scan(
			&timeBasedCode.Id,
			&timeBasedCode.Email,
			&timeBasedCode.Website,
			&timeBasedCode.Secret,
			&timeBasedCode.Type,
			&timeBasedCode.Validity,
			&timeBasedCode.CreatedAt,
			&timeBasedCode.UpdatedAt,
		)

		if err != nil {
			return []models.TimeBasedCode{}, err
		}

		timeBasedCodes = append(timeBasedCodes, timeBasedCode)
	}

	return timeBasedCodes, nil
}

func (tr *TotpRepository) Get(id uuid.UUID) (models.TimeBasedCode, error) {
	query := `
	SELECT t.id, t.website, t.email, LENGTH(t.secret) as password_lenght, tt.name, t.expiration, t.created_at, t.updated_at
	FROM public.time_based_codes as t
	JOIN public.time_based_code_types as tt
	on tt.id = t.type_id
	WHERE t.id=$1
	`
	queryResult := tr.Storage.Single(query, []interface{}{
		&id,
	})
	var timeBasedCode models.TimeBasedCode
	err := queryResult.Scan(
		&timeBasedCode.Id,
		&timeBasedCode.Email,
		&timeBasedCode.Website,
		&timeBasedCode.Secret,
		&timeBasedCode.Type,
		&timeBasedCode.Validity,
		&timeBasedCode.CreatedAt,
		&timeBasedCode.UpdatedAt,
	)

	if err != nil {
		return models.TimeBasedCode{}, err
	}

	return timeBasedCode, nil

}

func (tr *TotpRepository) GetCodeTypes() ([]models.CodeType, error) {
	query := `
		SELECT * FROM public.time_based_code_types
	`
	queryResult := tr.Storage.Where(query, []interface{}{})

	var codeTypes []models.CodeType
	for queryResult.Next() {
		var codeType models.CodeType
		err := queryResult.Scan(
			&codeType.Id,
			&codeType.Name,
		)

		if err != nil {
			return []models.CodeType{}, err
		}

		codeTypes = append(codeTypes, codeType)
	}

	return codeTypes, nil
}

func (sr *TotpRepository) Content(id uuid.UUID) (string, error) {
	query := `
		SELECT secret
		FROM public.time_based_codes
		WHERE id = $1
	`

	result := sr.Storage.Single(query, []interface{}{
		&id,
	})

	var codeSecret string
	err := result.Scan(&codeSecret)

	if err != nil {
		return "", err
	}

	return codeSecret, nil
}

func (tr *TotpRepository) Add(timePassword dtos.CreateTimeBasedCode) (uuid.UUID, error) {
	query := `
	INSERT INTO public.time_based_codes(
		website, email, secret, type_id, expiration, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	queryResult := tr.Storage.Single(query, []interface{}{
		&timePassword.Website,
		&timePassword.Email,
		&timePassword.Secret,
		&timePassword.Type,
		&timePassword.Validity,
		time.Now().UTC(),
	})

	var id uuid.UUID

	err := queryResult.Scan(&id)

	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (tr *TotpRepository) Update(timePassword *dtos.UpdateTimeBasedCode) bool {
	query := `
		UPDATE public.time_based_codes
		SET  website=$1, email=$2, secret=$3, type_id=$4, expiration=$5, updated_at=$6
		WHERE id= $7
	`

	result := tr.Storage.Exec(query, []interface{}{
		&timePassword.Website,
		&timePassword.Email,
		&timePassword.Secret,
		&timePassword.Type,
		&timePassword.Validity,
		time.Now().UTC(),
		&timePassword.Id,
	})

	return result
}

func (tr *TotpRepository) Remove(id uuid.UUID) bool {
	query := `
	DELETE FROM public.time_based_codes
	WHERE id = $1
`
	result := tr.Storage.Exec(query, []interface{}{
		&id,
	})

	return result
}
