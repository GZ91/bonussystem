package config

type Config struct {
	addressPort     string
	addressBaseData string
	addressAccrual  string
	secretKey       string
}

func New(addressPort, addressBaseData, addressAccrual, secretKey string) *Config {
	return &Config{
		addressPort:     addressPort,
		addressBaseData: addressBaseData,
		addressAccrual:  addressAccrual,
		secretKey:       secretKey,
	}
}

func (r *Config) GetAddressPort() string {
	return r.addressPort
}

func (r *Config) GetAddressBaseData() string {
	return r.addressBaseData
}

func (r *Config) GetAddressAccrual() string {
	return r.addressAccrual
}

func (r *Config) GetSecretKey() string {
	return r.secretKey
}
