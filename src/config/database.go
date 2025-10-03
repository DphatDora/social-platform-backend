package config

type Database struct {
	Username               string
	Password               string
	Host                   string
	Port                   string
	Name                   string
	Debug                  bool
	PoolSize               int
	IdleConnTimeoutSeconds int
	MaxConnAgeSeconds      int
	SslMode                string
	TimeZone               string
}
