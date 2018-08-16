package api

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

//Settings stores some configs about how the API will woks.
type Settings struct {
	Host          string
	Authorize     bool
	Authorization struct {
		JwksURI  string
		Issuer   string
		Audience []string
	}
	BasePath        string
	ApplicationName string
}

//SettingGenerator creates a instance of Settings.
type SettingGenerator func() *Settings

//DotEnv generates settings using environment variables.
func DotEnv(files ...string) SettingGenerator {
	return func() *Settings {
		err := godotenv.Load(files...)

		if err != nil {
			panic(err)
		}

		authorize, err := strconv.ParseBool(os.Getenv("AUTHORIZE"))

		if err != nil {
			authorize = true
		}

		authorization := struct {
			JwksURI  string
			Issuer   string
			Audience []string
		}{
			Audience: strings.Split(os.Getenv("AUDIENCE"), ","),
			Issuer:   os.Getenv("ISSUSER"),
			JwksURI:  os.Getenv("JWKS_URI"),
		}

		host := os.Getenv("HOST")
		basePath := os.Getenv("BASE_PATH")
		appName := os.Getenv("APPLICATION_NAME")

		return &Settings{
			Host:            host,
			BasePath:        basePath,
			Authorize:       authorize,
			Authorization:   authorization,
			ApplicationName: appName,
		}
	}
}
