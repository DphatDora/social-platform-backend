package config

type Auth struct {
	VerifyTokenExpirationMinutes int
	AccessTokenExpirationMinutes int
	JWTSecret                    string
}
