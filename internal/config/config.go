// Package config gère la configuration de l'agent
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config contient la configuration de l'agent
type Config struct {
	// Port d'écoute du serveur WebSocket
	Port int `json:"port"`

	// Origines autorisées (CORS)
	AllowedOrigins []string `json:"allowed_origins"`

	// Imprimante par défaut
	DefaultPrinter string `json:"default_printer,omitempty"`
}

// DefaultConfig retourne la configuration par défaut
func DefaultConfig() *Config {
	return &Config{
		Port: 19195,
		AllowedOrigins: []string{
			"https://billettera.com",
			"https://www.billettera.com",
			"http://localhost:8000",  // Dev Laravel
			"http://127.0.0.1:8000",  // Dev Laravel
		},
	}
}

// Load charge la configuration depuis un fichier JSON
// Si le fichier n'existe pas, retourne la configuration par défaut
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Chercher config.json à côté de l'exécutable
	exePath, err := os.Executable()
	if err != nil {
		return cfg, nil
	}

	configPath := filepath.Join(filepath.Dir(exePath), "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Fichier non trouvé, utiliser la config par défaut
		return cfg, nil
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// IsOriginAllowed vérifie si une origine est autorisée
func (c *Config) IsOriginAllowed(origin string) bool {
	for _, allowed := range c.AllowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}
