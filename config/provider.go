package config

type EnvProvider struct {
	app *App
	aws *AwsConfig
}

func (c EnvProvider) Aws() *AwsConfig {
	return c.aws
}

func (c EnvProvider) App() *App {
	return c.app
}

func NewEnvProvider() *EnvProvider {
	result := &EnvProvider{
		app: newAppConfig(),
		aws: newAWSConfig(),
	}

	return result
}
