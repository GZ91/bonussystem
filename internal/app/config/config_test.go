package config

import "testing"

func TestConfig(t *testing.T) {
	addressPort := "localhost:8080"
	addressBaseData := "http://api.example.com"
	addressAccrual := "http://accrual.example.com"
	secretKey := "mySecretKey"

	config := New(addressPort, addressBaseData, addressAccrual, secretKey)

	// Проверка значений, полученных через геттеры
	if config.GetAddressPort() != addressPort {
		t.Errorf("Expected AddressPort to be %s, got %s", addressPort, config.GetAddressPort())
	}

	if config.GetAddressBaseData() != addressBaseData {
		t.Errorf("Expected AddressBaseData to be %s, got %s", addressBaseData, config.GetAddressBaseData())
	}

	if config.GetAddressAccrual() != addressAccrual {
		t.Errorf("Expected AddressAccrual to be %s, got %s", addressAccrual, config.GetAddressAccrual())
	}

	if config.GetSecretKey() != secretKey {
		t.Errorf("Expected SecretKey to be %s, got %s", secretKey, config.GetSecretKey())
	}
}
