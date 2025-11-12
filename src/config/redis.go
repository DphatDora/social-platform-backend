package config

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
	Required bool
}
