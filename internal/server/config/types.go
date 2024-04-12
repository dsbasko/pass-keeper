package config

import "time"

type config struct {
	Env        string `env:"ENV" default:"prod"`
	AppName    string `env:"APP_NAME"`
	Controller controller
	Provider   provider
}

type controller struct {
	GRPC controllerGRPC
}

type controllerGRPC struct {
	Port      int `env:"GRPC_PORT"`
	TimeoutMs int `env:"GRPC_TIMEOUT_MS"`
}

type provider struct {
	Postgre postgre
}

type postgre struct {
	Host     string `env:"POSTGRE_HOST"`
	Port     int    `env:"POSTGRE_PORT"`
	User     string `env:"POSTGRE_USER"`
	Pass     string `env:"POSTGRE_PASS"`
	DB       string `env:"POSTGRE_DB"`
	MaxConns int    `env:"POSTGRE_MAX_CONNS"`
	DSN      string `env:"POSTGRE_DSN"`
}

func (c *controllerGRPC) Timeout() time.Duration {
	return time.Duration(c.TimeoutMs) * time.Millisecond
}
