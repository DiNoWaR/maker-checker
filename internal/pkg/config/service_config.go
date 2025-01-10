package config

type ServiceConfig struct {
	ServiceHost string `env:"SERVICE_HOST"`
	ServicePort string `env:"SERVICE_PORT"`
	DBConfig    DBConfig
}

type DBConfig struct {
	Host     string `env:"SERVICE_DB_HOST"`
	Port     string `env:"SERVICE_DB_PORT"`
	Database string `env:"SERVICE_DB_NAME"`
	Username string `env:"SERVICE_DB_USER"`
	Password string `env:"SERVICE_DB_PASSWORD"`
}
