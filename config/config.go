package config

import (
	"os"
)

//function to take port value from .env file
func ConfigPort() string {

	return os.Getenv("PORT")
}

//function to take mongodb url value from .env file
func ConfigDB() string {

	return os.Getenv("DATABASE_URL")
}

func ConfigDBname() string {

	return os.Getenv("DATABASE_NAME")
}

//function to take openai api key value from .env file
func ConfigOpenAI() string {

	return os.Getenv("OPENAI_API_KEY")
}

//function to take jwt key value from .env file
func ConfigJWTkey() string {

	return os.Getenv("JWT_KEY")
}