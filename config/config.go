package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"runtime"
	"time"
)

type Config struct {
	APP struct {
		TeamsList    []string      `env:"APP_TEAMS_LIST,default=Gremio"`
		CacheEnabled bool          `env:"CACHE_ENABLED,default=true"`
		CachePath    string        `env:"CACHE_PATH,default=$HOME/.acmeinc/another-go-challenge/bruno-nascimento"`
		CacheTTL     time.Duration `env:"CACHE_TTL,default=1s"`
	}
	HTTP struct {
		Port    string        `env:"HTTP_PORT,default=:8080"`
		Timeout time.Duration `env:"HTTP_TIMEOUT,default=3s"`
	}
	API struct {
		Endpoint         string        `env:"API_ENDPOINT,default=https://acme.inc/some-path/{team_id}.json"`
		ParallelRequests int           `env:"API_PARALLEL_REQUESTS"`
		RequestTimeout   time.Duration `env:"API_REQUEST_TIMEOUT,default=3s"`
	}
}

func New() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process(context.Background(), config); err != nil {
		return nil, err
	}
	setParallelRequests(config)
	return config, nil
}

func NewMock(mapper map[string]string) (*Config, error) {
	lookups := envconfig.MapLookuper(mapper)
	config := &Config{}
	err := envconfig.ProcessWith(context.Background(), config, lookups)
	if err != nil {
		return nil, err
	}
	setParallelRequests(config)
	return config, nil
}

func setParallelRequests(cfg *Config) {
	if cfg.API.ParallelRequests == 0 {
		cfg.API.ParallelRequests = runtime.NumCPU() * 4
	}
}
