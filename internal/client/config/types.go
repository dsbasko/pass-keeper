package config

type config struct {
	Env     string `env:"ENV" default:"prod"`
	AppName string `env:"APP_NAME"`
	Server  server
}

type server struct {
	GRPC serverGRPC
}

type serverGRPC struct {
	Port int `env:"SERVER_GRPC_PORT"`
}
