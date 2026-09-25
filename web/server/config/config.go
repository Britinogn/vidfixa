package config

import "os"

type Config struct {
	AppName   string
	AppEnv    string
	AppPort   string
	AppURL    string
	CorsURL string

	DatabaseURL string

	JWTSecret    string
	JWTExpiresIn string

	DownloadWorkers string
	YtdlpPath       string
	FFmpegPath      string

	BachsAPIKey        string
	BachsBaseURL       string
	BachsWebhookSecret string
	BachsPlusProductID string
	BachsProProductID  string

	BACHS_SUCCESS_URL string
	BACHS_CANCEL_URL  string

	BcryptRounds string
	DownloadDir  string
}

func Load() *Config {
	return &Config{
		AppName:   os.Getenv("APP_NAME"),
		AppEnv:    os.Getenv("APP_ENV"),
		AppPort:   os.Getenv("APP_PORT"),
		AppURL:    os.Getenv("APP_URL"),
		CorsURL: os.Getenv("CORS_URL"),

		DatabaseURL: os.Getenv("DATABASE_URL"),

		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: os.Getenv("JWT_EXPIRES_IN"),

		DownloadWorkers: os.Getenv("DOWNLOAD_WORKERS"),
		YtdlpPath:       os.Getenv("YTDLP_PATH"),
		FFmpegPath:      os.Getenv("FFMPEG_PATH"),

		BachsAPIKey:        os.Getenv("BACHS_API_KEY"),
		BachsBaseURL:       os.Getenv("BACHS_BASE_URL"),
		BachsPlusProductID: os.Getenv("BACHS_PLUS_PRODUCT_ID"),
		BachsWebhookSecret: os.Getenv("BACHS_WEBHOOK_SECRET"),
		BachsProProductID:  os.Getenv("BACHS_PRO_PRODUCT_ID"),
		DownloadDir:        os.Getenv("DOWNLOAD_DIR"),

		// sample
		BACHS_SUCCESS_URL: os.Getenv("BACHS_SUCCESS_URL"),
		BACHS_CANCEL_URL:  os.Getenv("BACHS_CANCEL_URL"),

		BcryptRounds: os.Getenv("BCRYPT_ROUNDS"),
	}
}
