package config

import (
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var (
	config *Config
	once sync.Once
)

type Config struct {
	MongoURI string
	MongoDatabase string

	HTTPPort string
	CertFile, KeyFile string

	Debug bool
}

func Get(path string) *Config {
	once.Do(func() {
		_ = godotenv.Load(path)
		config = &Config{
			MongoURI:      os.Getenv("MONGO_URI"),
			MongoDatabase: os.Getenv("MONGO_DATABASE"),
			HTTPPort:      os.Getenv("HTTP_PORT"),
			Debug: 		   os.Getenv("DEBUG") == "true",
		}
	})

	return config
}

