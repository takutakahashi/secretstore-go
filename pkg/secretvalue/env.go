package secretvalue

import "encoding/json"

type EnvSecret struct {
	envVars map[string]string
}

func NewEnvSecret(envVars map[string]string) *EnvSecret {
	return &EnvSecret{
		envVars: envVars,
	}
}

func (s EnvSecret) Data() ([]byte, error) {
	return json.Marshal(s.envVars)
}
