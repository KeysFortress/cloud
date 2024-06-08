package models

type CredentialsData struct {
	Passwords                 int
	LastUsedPasswords         int
	Secrets                   int
	LastUsedSecrets           int
	Identities                int
	LastUsedIdentities        int
	TimeBasedCodes            int
	LastUsedTimeBasedCodes    int
	MonthlyCredentialRequests int
	NewCredentials            int
}
