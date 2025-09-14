package config

type App struct {
	Host         string
	Port         int
	Debug        bool
	Whitelist    []string
	ReadTimeout  int
	WriteTimeout int
}
