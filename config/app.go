package config

import (
	"os"
)

func newAppConfig() *App {
	return &App{
		app:        os.Getenv("APP"),
		version:    os.Getenv("VERSION"),
		env:        os.Getenv("ENV"),
		smsHost:    os.Getenv("SMS_HOST"),
		smsAuthKey: os.Getenv("SMS_AUTH_KEY"),
	}
}

type App struct {
	app     string
	version string
	env     string

	smsHost    string
	smsAuthKey string
}

func (a App) SmsHost() string {
	return a.smsHost
}

func (a App) SmsAuthKey() string {
	return a.smsAuthKey
}

func (a App) App() string {
	return a.app
}

func (a App) Version() string {
	return a.version
}

func (a App) Env() string {
	return a.env
}
