package initializing

import (
	"github.com/GZ91/bonussystem/internal/app/config"
	"github.com/GZ91/bonussystem/internal/app/initializing/envs"
	"github.com/GZ91/bonussystem/internal/app/initializing/flags"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"go.uber.org/zap"
)

func Configuration() *config.Config {
	logger.Initializing("info")
	conf := config.New(ReadParams())
	return conf
}

func ReadParams() (string, string, string, string) {
	envVars, err := envs.ReadEnv()
	if err != nil {
		logger.Log.Error("error when reading environment variables", zap.Error(err))
	}
	var AddressPort, AddressBaseData, AddressAccrual string
	if envVars == nil {
		AddressPort, AddressBaseData, AddressAccrual = flags.ReadFlags()
	} else {
		AddressPort, AddressBaseData, AddressAccrual = envVars.AddressPort, envVars.AddressBaseData, envVars.AddressAccrual

		if AddressPort == "" || AddressBaseData == "" || AddressAccrual == "" {
			AddressPortFlag, AddressBaseDataFlag, AddressAccrualFlag := flags.ReadFlags()
			if AddressPort == "" {
				AddressPort = AddressPortFlag
			}
			if AddressBaseData == "" {
				AddressBaseData = AddressBaseDataFlag
			}
			if AddressAccrual == "" {
				AddressAccrual = AddressAccrualFlag
			}
		}
	}
	SecretKey := envVars.SecretKey
	if SecretKey == "" {
		SecretKey = "SecretKey"
	}
	return AddressPort, AddressBaseData, AddressAccrual, SecretKey
}
