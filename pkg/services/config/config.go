package config

import "os"

type postgresConfig struct {
	PostgresUser     string
	PostgresPassword string
	PostgresAddress  string
	PostgresPort     string
	PostgresDB       string
}

type redisConfig struct {
	RedisURL string
}

type sessionConfig struct {
	JWTKey     string
	SessionKey string
}

type grpcConfig struct {
	GRPCPort string
}

type Config struct {
	RedisConfig    redisConfig
	PostgresConfig postgresConfig
	SessionConfig  sessionConfig
	GRPCConfig     grpcConfig
}

func NewConfig() *Config {
	return &Config{
		RedisConfig: redisConfig{
			getEnv("REDIS_URL", ""),
		},
		PostgresConfig: postgresConfig{
			PostgresUser:     getEnv("POSTGRES_USER", ""),
			PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
			PostgresAddress:  getEnv("POSTGRES_ADDRESS", ""),
			PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
			PostgresDB:       getEnv("POSTGRES_DATABASE", ""),
		},
		SessionConfig: sessionConfig{
			JWTKey:     getEnv("JWT_SECRET_KEY", ""),
			SessionKey: getEnv("SESSION_KEY", ""),
		},
		GRPCConfig: grpcConfig{
			GRPCPort: getEnv("GRPC_PORT", ""),
		},
	}
}

func getEnv(name string, defaultValue string) string {
	if value, isExists := os.LookupEnv(name); isExists {
		return value
	}
	return defaultValue
}
