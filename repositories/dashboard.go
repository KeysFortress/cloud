package repositories

import (
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type DashboardRepository struct {
	Storage interfaces.Storage
}

func (d *DashboardRepository) Credentials() (models.CredentialsData, error) {
	sql := `
	WITH MonthlyCredentialRequests AS (
		SELECT
			COUNT(*) AS RequestCount
		FROM (
			SELECT id FROM account_identities WHERE EXTRACT(MONTH FROM updated_at) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM updated_at) = EXTRACT(YEAR FROM CURRENT_DATE)
			UNION ALL
			SELECT id FROM account_passwords WHERE EXTRACT(MONTH FROM updated_at) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM updated_at) = EXTRACT(YEAR FROM CURRENT_DATE)
			UNION ALL
			SELECT id FROM account_secrets WHERE EXTRACT(MONTH FROM updated_at) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM updated_at) = EXTRACT(YEAR FROM CURRENT_DATE)
			UNION ALL
			SELECT id FROM time_based_codes WHERE EXTRACT(MONTH FROM updated_at) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM updated_at) = EXTRACT(YEAR FROM CURRENT_DATE)
		) AS all_requests
	),
	NewCredentials AS (
		SELECT
			COUNT(*) AS NewCount
		FROM (
			SELECT id FROM account_identities WHERE DATE_TRUNC('week', created_at) = DATE_TRUNC('week', CURRENT_DATE)
			UNION ALL
			SELECT id FROM account_passwords WHERE DATE_TRUNC('week', created_at) = DATE_TRUNC('week', CURRENT_DATE)
			UNION ALL
			SELECT id FROM account_secrets WHERE DATE_TRUNC('week', created_at) = DATE_TRUNC('week', CURRENT_DATE)
			UNION ALL
			SELECT id FROM time_based_codes WHERE DATE_TRUNC('week', created_at) = DATE_TRUNC('week', CURRENT_DATE)
		) AS new_creds
	)
	SELECT
		(SELECT COUNT(*) FROM account_passwords) AS Passwords,
		(SELECT COUNT(*) FROM account_passwords WHERE updated_at IS NOT NULL) AS LastUsedPasswords,
		(SELECT COUNT(*) FROM account_secrets) AS Secrets,
		(SELECT COUNT(*) FROM account_secrets WHERE updated_at IS NOT NULL) AS LastUsedSecrets,
		(SELECT COUNT(*) FROM account_identities) AS Identities,
		(SELECT COUNT(*) FROM account_identities WHERE updated_at IS NOT NULL) AS LastUsedIdentities,
		(SELECT COUNT(*) FROM time_based_codes) AS TimeBasedCodes,
		(SELECT COUNT(*) FROM time_based_codes WHERE updated_at IS NOT NULL) AS LastUsedTimeBasedCodes,
		(SELECT RequestCount FROM MonthlyCredentialRequests) AS MonthlyCredentialRequests,
		(SELECT NewCount FROM NewCredentials) AS NewCredentials;
	`

	query := d.Storage.Single(sql, []interface{}{})

	var credentialRequest models.CredentialsData
	err := query.Scan(
		&credentialRequest.Passwords,
		&credentialRequest.LastUsedPasswords,
		&credentialRequest.Secrets,
		&credentialRequest.LastUsedSecrets,
		&credentialRequest.Identities,
		&credentialRequest.LastUsedIdentities,
		&credentialRequest.TimeBasedCodes,
		&credentialRequest.LastUsedTimeBasedCodes,
		&credentialRequest.MonthlyCredentialRequests,
		&credentialRequest.NewCredentials,
	)

	if err != nil {
		return models.CredentialsData{}, err
	}

	return credentialRequest, nil
}

func (d *DashboardRepository) Devices() ([]models.Device, error) {
	sql := `
		SELECT aaak.name, aaak.last_used_at, dt.name, dt.id from associated_account_access_keys as aaak
		JOIN device_types as dt on  dt.id = aaak.device_type_id
	`

	query := d.Storage.Where(sql, []interface{}{})

	var devices []models.Device
	for query.Next() {
		var device models.Device

		err := query.Scan(&device.Name, &device.LastUsed, &device.Type, &device.TypeId)

		if err != nil {
			return []models.Device{}, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}
