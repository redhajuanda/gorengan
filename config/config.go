package config

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// Config represents configuration variables
type Config struct {
	Server struct {
		ENV  string `envconfig:"APP_ENV"`
		PORT string `envconfig:"APP_PORT"`
	}
	JWT struct {
		SigningKey      string `envconfig:"JWT_SIGNING_KEY"`
		TokenExpiration int    `envconfig:"JWT_TOKEN_EXPIRATION"`
	}
	Database struct {
		Host     string `envconfig:"DB_HOST"`
		Port     string `envconfig:"DB_PORT"`
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
		DBName   string `envconfig:"DB_NAME"`
	}
}

// LoadTest loads test config
func LoadTest() Config {
	return load("test", ".env.test")
}

// LoadDefault loads default config (default.yml) and override config with env if supplied
func LoadDefault() Config {
	return load("default", ".env")
}

// load config and populate to config struct
// load will load config supplied file without extension,
// and will override config with env if env supplied
func load(file string, env string) Config {
	path := getSourcePath()
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error load config file: %v", err)
	}
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshal config: %v", err)
	}
	if env != "" {
		readEnv(&config, env) // Read and OVERWRITE yml with env
	}

	return config
}

func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func readEnv(cfg *Config, env string) {
	err := godotenv.Load(getSourcePath() + "/../" + env)
	if err != nil {
		fmt.Println(err)
	}
	err = envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("Error populate env: %v", err)
	}
}
