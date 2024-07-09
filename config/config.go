package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Env        string `env:"ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"TODO_DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"TODO_DB_PORT" envDefault:"33306"`
	DBUser     string `env:"TODO_DB_USER" envDefault:"todo"`
	DBPassword string `env:"TODO_DB_PASSWORD" envDefault:"todo"`
	DBName     string `env:"TODO_DB_name" envDefault:"todo"`
	RedisHost  string `env:"TODO_REDIS_HOST" envDefault:"localhost"`
	RedisPort  int    `env:"TODO_REDIS_PORT" envDefault:"36379"`
}

func New() (*Config, error) {
	c := &Config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	return c, nil
}
