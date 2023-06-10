package config

type Config struct {
	addressPort     string
	addressBaseData string
	addressAccrual  string
}

func New(addressPort, addressBaseData, addressAccrual string) *Config {
	return &Config{
		addressPort:     addressPort,
		addressBaseData: addressBaseData,
		addressAccrual:  addressAccrual,
	}
}
