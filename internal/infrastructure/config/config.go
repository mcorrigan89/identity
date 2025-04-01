package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN     string
		Logging bool
	}
	CientURL    string
	ServerToken string
	Cors        struct {
		TrustedOrigins []string
	}
	Storage struct {
		Endpoint  string
		AccessKey string
		SecretKey string
	}
	OAuth struct {
		Google struct {
			ClientID     string
			ClientSecret string
		}
	}
	Mail struct {
		SMTPServer   string
		SMTPUsername string
		SMTPPassword string
	}
}

func LoadConfig(cfg *Config) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Load ENV
	env := os.Getenv("ENV")
	if env == "" {
		cfg.Env = "local"
	} else {
		cfg.Env = env
	}

	// Load PORT
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("PORT not available in .env")
	}

	cfg.Port = port

	// Load CLIENT_URL
	client_url := os.Getenv("CLIENT_URL")
	if client_url == "" {
		log.Fatalf("CLIENT_URL not available in .env")
	}

	cfg.CientURL = client_url

	// Load SERVER_TOKEN
	server_token := os.Getenv("SERVER_TOKEN")
	if server_token == "" {
		log.Fatalf("SERVER_TOKEN not available in .env")
	}

	cfg.ServerToken = server_token

	// Load DATABASE_URL
	postgres_url := os.Getenv("POSTGRES_URL")
	if postgres_url == "" {
		log.Fatalf("POSTGRES_URL not available in .env")
	}

	cfg.DB.DSN = postgres_url

	cfg.Cors.TrustedOrigins = []string{"http://localhost:3000"}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		log.Fatalf("MINIO_ENDPOINT not available in .env")
	}
	cfg.Storage.Endpoint = endpoint

	access_key := os.Getenv("MINIO_ACCESS_KEY")
	if access_key == "" {
		log.Fatalf("MINIO_ACCESS_KEY not available in .env")
	}
	cfg.Storage.AccessKey = access_key

	secret_key := os.Getenv("MINIO_SECRET_KEY")
	if secret_key == "" {
		log.Fatalf("MINIO_SECRET_KEY not available in .env")
	}
	cfg.Storage.SecretKey = secret_key

	smtp_server := os.Getenv("SMTP_SERVER")
	if smtp_server == "" {
		log.Fatalf("SMTP_SERVER not available in .env")
	}
	cfg.Mail.SMTPServer = smtp_server

	smtp_username := os.Getenv("SMTP_USERNAME")
	if smtp_username == "" {
		log.Fatalf("SMTP_USERNAME not available in .env")
	}
	cfg.Mail.SMTPUsername = smtp_username

	smtp_password := os.Getenv("SMTP_PASSWORD")
	if smtp_password == "" {
		log.Fatalf("SMTP_PASSWORD not available in .env")
	}
	cfg.Mail.SMTPPassword = smtp_password

	google_client_id := os.Getenv("GOOGLE_CLIENT_ID")
	if google_client_id == "" {
		log.Fatalf("GOOGLE_CLIENT_ID not available in .env")
	}

	google_client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if google_client_secret == "" {
		log.Fatalf("GOOGLE_CLIENT_SECRET not available in .env")
	}

	cfg.OAuth.Google.ClientID = google_client_id
	cfg.OAuth.Google.ClientSecret = google_client_secret
}
