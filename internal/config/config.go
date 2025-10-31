package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleOAuthConfig *oauth2.Config
	SessionSecret     string
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	// Carrega o .env
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env")
	}

	AppConfig = &Config{
		SessionSecret: os.Getenv("SESSION_SECRET"),
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}

	if AppConfig.SessionSecret == "" {
		log.Fatal("SESSION_SECRET não está definida no .env")
	}

	return AppConfig, nil
}
