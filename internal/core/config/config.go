package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v10"
)

const (
	// Test environment.
	Test Environment = "TEST"
	// Production environment.
	Production Environment = "PROD"
)

type (
	// Environment type to hold values from it.
	Environment string

	// Configuration contains loaded environment variables.
	Configuration struct {
		Environment Environment `env:"ENVIRONMENT,required"`
		LogLevel    int         `env:"LOG_LEVEL" envDefault:"-1"` // Debug = -1
		Application AppMetadata `envPrefix:"APP_"`
		Server      Server      `envPrefix:"SERVER_"`
		Ingestor    Ingestor    `envPrefix:"INGESTOR_"`
	}

	// AppMetadata contains the application's metadata.
	AppMetadata struct {
		Name    string `env:"NAME,required"`
		Version string `env:"VERSION,required"`
	}

	// Server contains server environment variables.
	Server struct {
		Port int `env:"PORT" envDefault:"8080"`
	}

	// Ingestor contains ingestor environment variables.
	Ingestor struct {
		Filepath string `env:"FILEPATH"`
	}
)

// Load loads values from environment variables into the Configuration struct.
func Load() (Configuration, error) {
	var c Configuration

	if err := env.Parse(&c); err != nil {
		return Configuration{}, fmt.Errorf("failed to load config: %w", err)
	}

	return c, nil
}

func (e Environment) IsProduction() bool {
	return strings.ToUpper(string(e)) == string(Production)
}
