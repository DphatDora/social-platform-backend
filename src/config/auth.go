package config

type Auth struct {
	VerifyTokenExpirationMinutes int
	AccessTokenExpirationMinutes int
	ResetTokenExpirationMinutes  int
	JWTSecret                    string
	GoogleClientID               string
	GoogleClientSecret           string
	GoogleRedirectURL            string
}
