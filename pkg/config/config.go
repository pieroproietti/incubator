package config

import (
	"os"
	"strconv"
)

// Config contiene tutte le variabili d'ambiente passate dalla CI/CD
type Config struct {
	TargetDir  string
	BaseVMID   int
	Firmware   string
	FsType     string
	Storage    string
	IsoStorage string
	Template   string
	Bridge     string
}

// GetEnv estrae una stringa o usa il default
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetEnvAsInt estrae un intero o usa il default
func GetEnvAsInt(key string, fallback int) int {
	strValue := GetEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}
