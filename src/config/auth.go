package config

type Auth struct {
	VerifyTokenExpirationMinutes int
	AccessTokenExpirationMinutes int
	RefreshTokenExpirationDays   int
	ResetTokenExpirationMinutes  int
	JWTSecret                    string
	GoogleClientID               string
	GoogleClientSecret           string
	GoogleRedirectURL            string
}
