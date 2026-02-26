package config

import (
	"os"
	"strconv"
	"strings"
)

func newAppConfig() *App {
	mailTrapPort, _ := strconv.Atoi(os.Getenv("MAIL_TRAP_PORT"))

	return &App{
		app:              os.Getenv("APP"),
		version:          os.Getenv("VERSION"),
		env:              os.Getenv("ENV"),
		emailFrom:        os.Getenv("EMAIL_FROM"),
		mailTransport:    os.Getenv("MAIL_TRANSPORT"),
		smsHost:          os.Getenv("SMS_HOST"),
		smsAuthKey:       os.Getenv("SMS_AUTH_KEY"),
		mailTrapHost:     os.Getenv("MAIL_TRAP_HOST"),
		mailTrapPort:     mailTrapPort,
		mailTrapUser:     os.Getenv("MAIL_TRAP_USER"),
		mailTrapPassword: os.Getenv("MAIL_TRAP_PASSWORD"),
		blackListEmails:  strings.Split(os.Getenv("BLACK_LIST_EMAILS"), ","),
	}
}

type App struct {
	app     string
	version string
	env     string

	emailFrom     string
	mailTransport string

	smsHost    string
	smsAuthKey string

	mailTrapHost     string
	mailTrapPort     int
	mailTrapUser     string
	mailTrapPassword string
	blackListEmails  []string
}

func (a App) EmailFrom() string {
	return a.emailFrom
}

func (a App) SmsHost() string {
	return a.smsHost
}

func (a App) MailTrapPassword() string {
	return a.mailTrapPassword
}

func (a App) BlackListEmails() []string {
	return a.blackListEmails
}

func (a App) MailTrapPort() int {
	return a.mailTrapPort
}

func (a App) SmsAuthKey() string {
	return a.smsAuthKey
}

func (a App) MailTransport() string {
	return a.mailTransport
}

func (a App) MailTrapHost() string {
	return a.mailTrapHost
}

func (a App) MailTrapUser() string {
	return a.mailTrapUser
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
