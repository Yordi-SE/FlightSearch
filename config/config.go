package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	ClientID     string
	ClientSecret string
	PCC          string
	URL          string
	SABREAUTHURL string
}

// New returns a new Config instance
func New() (*Config, error) {
	c := &Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
		PCC:          os.Getenv("PCC"),
		URL:          os.Getenv("URL"),
		SABREAUTHURL: os.Getenv("SABREAUTHURL"),
	}

	if c.ClientID == "" {
		return nil, fmt.Errorf("CLIENT_ID is required")
	}
	if c.ClientSecret == "" {
		return nil, fmt.Errorf("CLIENT_SECRET is required")
	}
	if c.PCC == "" {
		return nil, fmt.Errorf("PCC is required")
	}
	if c.SABREAUTHURL == "" {
		return nil, fmt.Errorf("SABREAUTHURL is required")
	}
	if c.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	return c, nil
}
