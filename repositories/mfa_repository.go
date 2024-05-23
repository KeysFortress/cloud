package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type MfaRepository struct {
	Storage interfaces.Storage
}

func (m *MfaRepository) IsConfigured(id uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM mfa_methods
		WHERE user_id = $1
	`

	queryResult := m.Storage.Single(query, []interface{}{

		&id,
	})

	var count int
	err := queryResult.Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (m *MfaRepository) GetForUser(id uuid.UUID) ([]models.UserMfa, error) {
	query := `
		SELECT id, type_id, value, address FROM mfa_methods
		WHERE user_id = $1
	`

	queryResult := m.Storage.Where(query, []interface{}{
		&id,
	})

	var result []models.UserMfa
	for queryResult.Next() {
		var mfaResult models.UserMfa

		err := queryResult.Scan(
			mfaResult.Id,
			mfaResult.TypeId,
			mfaResult.Value,
			mfaResult.Address,
		)
		if err != nil {
			return []models.UserMfa{}, err
		}

		result = append(result, mfaResult)
	}

	return result, nil
}

func (m *MfaRepository) GetForUserByType(id uuid.UUID, typeId int) ([]models.UserMfa, error) {
	query := `
		SELECT id, type_id, address value FROM mfa_methods
		WHERE user_id = $1 and type_id = $2
	`

	queryResult := m.Storage.Where(query, []interface{}{
		&id,
		typeId,
	})

	var result []models.UserMfa
	for queryResult.Next() {
		var mfaResult models.UserMfa

		err := queryResult.Scan(
			mfaResult.Id,
			mfaResult.TypeId,
			mfaResult.Value,
			mfaResult.Address,
		)

		if err != nil {
			return []models.UserMfa{}, err
		}

		result = append(result, mfaResult)
	}

	return result, nil
}

func (m *MfaRepository) Add(secret string, typeId int, user uuid.UUID, email sql.NullString) (bool, error) {
	qeury := `
	INSERT INTO public.mfa_methods(
		type_id, value, user_id, address)
		VALUES ($2, $1, $3, $4)
		RETURNING id;
	`

	queryResult := m.Storage.Single(qeury, []interface{}{
		secret,
		typeId,
		user,
		email,
	})

	var recordId uuid.UUID
	err := queryResult.Scan(&recordId)

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func (m *MfaRepository) Delete(id uuid.UUID) bool {

	query := `
		DELETE FROM public.mfa_methods
		WHERE id = $1
	`
	result := m.Storage.Exec(query, []interface{}{
		&id,
	})

	return result
}
