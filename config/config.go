package config

import (
	"os"

	"github.com/joho/godotenv"
)

//function to take host value from .env file
func ConfigHost() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("HOST")
}

//function to take port value from .env file
func ConfigPort() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("PORT")
}

//function to take mongodb url value from .env file
func ConfigDB() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("DATABASE_URL")
}

//function to take openai api key value from .env file
func ConfigGeminiKey() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("GEMINI_API_KEY")
}

func ConfigGoogleOauthClientId() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
}

func ConfigGoogleOauthClientSecret() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
}

func ConfigGoogleOauthRedirectUrl() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
}

func ConfigFirebaseProjectId() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("FIREBASE_PROJECT_ID")
}