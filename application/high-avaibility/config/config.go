package config

type Config struct {
	EtcdEndpoints      string
	DatabaseConnection string
}

func NewConfig() *Config {
	return &Config{
		EtcdEndpoints:      "localhost:2379",
		DatabaseConnection: "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	}
}
