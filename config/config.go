package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string
}

type Config struct {
	DbConfig
	TokenConfig
}

type TokenConfig struct {
	IssuerName      string
	JwtSignatureKey []byte
	JwtLifeTime     time.Duration
}

func (c *Config) ReadConfig() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error Load .env file")
	}

	// Environment Variable
	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	token_lifetime, err := strconv.Atoi(os.Getenv("TOKEN_LIFE_TIME"))
	if err != nil {
		return errors.New("Fail Parse Token Life Time")
	}

	c.TokenConfig = TokenConfig{
		IssuerName:      os.Getenv("ISSUER_NAME"),
		JwtSignatureKey: []byte(os.Getenv("SIGNATURE")),
		JwtLifeTime:     time.Duration(token_lifetime) * time.Hour, //time.Minute,
	}

	// Cek jika tidak mengirimkan
	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" || c.DbConfig.User == "" || c.DbConfig.Password == "" || c.DbConfig.Driver == "" {
		return fmt.Errorf("Missing ENV")
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.ReadConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// SET DB_HOST=localhost
// SET ...

// go run .
