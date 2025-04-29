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

func (s EnvSecret) GetData() ([]byte, error) {
	return json.Marshal(s.envVars)
}

func (s *EnvSecret) SetData(data []byte) error {
	return json.Unmarshal(data, &s.envVars)
}
