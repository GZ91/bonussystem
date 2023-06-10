package envs

import "github.com/caarlos0/env/v6"

type EnvVars struct {
	AddressPort     string `env:"RUN_ADDRESS"`
	AddressBaseData string `env:"DATABASE_URI"`
	AddressAccrual  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func ReadEnv() (*EnvVars, error) {

	envs := EnvVars{}

	if err := env.Parse(&envs); err != nil {
		return nil, err
	}

	return &envs, nil
}
