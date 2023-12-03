package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		PG    `yaml:"postgres"`
		OAuth `yaml:"oauth"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PgUser string `env-required:"true" yaml:"user" env:"PG_USER"`
		PgPass string `env-required:"true" yaml:"password" env:"PG_PASS"`
		PgHost string `env-required:"true" yaml:"host" env:"PG_HOST"`
		PgName string `env-required:"true" yaml:"name" env:"PG_NAME"`
		PgPort string `env-required:"true" yaml:"port" env:"PG_PORT"`
	}

	// OAuth -.
	OAuth struct {
		GithubClientId     string `env-required:"true" yaml:"github_client_id" env:"GITHUB_CLIENT_ID"`
		GithubClientSecret string `env-required:"true" yaml:"github_client_secret" env:"GITHUB_CLIENT_SECRET"`
		GithubRedirectUri  string `env-required:"true" yaml:"github_redirect_uri" env:"GITHUB_REDIRECT_URI"`
		GoogleClientId     string `env-required:"true" yaml:"google_client_id" env:"GOOGLE_CLIENT_ID"`
		GoogleClientSecret string `env-required:"true" yaml:"google_client_secret" env:"GOOGLE_CLIENT_SECRET"`
		GoogleRedirectUri  string `env-required:"true" yaml:"google_redirect_uri" env:"GOOGLE_REDIRECT_URI"`
		RedditClientId     string `env-required:"true" yaml:"reddit_client_id" env:"REDDIT_CLIENT_ID"`
		RedditClientSecret string `env-required:"true" yaml:"reddit_client_secret" env:"REDDIT_CLIENT_SECRET"`
		RedditRedirectUri  string `env-required:"true" yaml:"reddit_redirect_uri" env:"REDDIT_REDIRECT_URI"`
		TokenEncryptSecret string `env-required:"true" yaml:"token_encrypt_secret" env:"TOKEN_ENCRYPT_SECRET"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
