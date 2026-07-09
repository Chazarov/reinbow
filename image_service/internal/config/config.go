// Пакет config загружает и предоставляет конфигурацию приложения.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config содержит параметры подключения к инфраструктуре и настройки сервиса.
type Config struct {
	DbHost                    string
	DbPort                    string
	DbUser                    string
	DbPassword                string
	DbName                    string
	DbSSLMode                 string
	RedisHost                 string
	RedisPort                 string
	ServerHost                string
	ServerPort                string
	ImageStoragePath          string
	NotifyServiceURL          string
	ExternalURL               string
	DefaultLogsPath           string
	CriticalLogsPath          string
	LoggerDebug               bool
	JwtSecret                 string
	JwtAccessTTL              time.Duration
	JwtRefreshTTL             time.Duration
	PaginationDefaultPage     int
	PaginationDefaultPageSize int
	TgAdminId                 int
}

// PaginationParams задаёт значения пагинации по умолчанию.
type PaginationParams struct {
	DefaultPage     int
	DefaultPageSize int
}

// PaginationConfig предоставляет параметры пагинации по умолчанию.
type PaginationConfig interface {
	GetPaginationParams() PaginationParams
}

// GetPaginationParams возвращает параметры пагинации по умолчанию из конфигурации.
func (c *Config) GetPaginationParams() PaginationParams {
	return PaginationParams{
		DefaultPage:     c.PaginationDefaultPage,
		DefaultPageSize: c.PaginationDefaultPageSize,
	}
}

// Dsn формирует строку подключения к PostgreSQL.
func (c *Config) Dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DbHost, c.DbPort, c.DbUser, c.DbPassword, c.DbName, c.DbSSLMode,
	)
}

// ServerAddr возвращает адрес HTTP-сервера в формате host:port.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

// RedisAddr возвращает адрес Redis в формате host:port.
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// LoadConfig загружает и валидирует конфигурацию приложения из переменных окружения.
func LoadConfig() (*Config, error) {
	required := map[string]*string{
		"DB_HOST":            nil,
		"DB_PORT":            nil,
		"DB_USER":            nil,
		"DB_PASSWORD":        nil,
		"DB_NAME":            nil,
		"DB_SSLMODE":         nil,
		"REDIS_HOST":         nil,
		"REDIS_PORT":         nil,
		"SERVER_HOST":        nil,
		"SERVER_PORT":        nil,
		"IMAGE_STORAGE_PATH": nil,
		"NOTIFY_SERVICE_URL": nil,
		"EXTERNAL_URL":       nil,
		"DEFAULT_LOGS_PATH":  nil,
		"CRITICAL_LOGS_PATH": nil,
		"JWT_SECRET":         nil,
		"JWT_ACCESS_TTL":     nil,
		"JWT_REFRESH_TTL":    nil,
		"TG_ADMIN_ID":        nil,
	}

	values := make(map[string]string, len(required))
	for key := range required {
		val := strings.TrimSpace(os.Getenv(key))
		if val == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
		values[key] = val
	}

	loggerDebug := false
	if raw := os.Getenv("LOGGER_DEBUG"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid LOGGER_DEBUG value %q: %w", raw, err)
		}
		loggerDebug = parsed
	}

	tgAdminId, err := strconv.Atoi(values["TG_ADMIN_ID"])
	if err != nil {
		return nil, fmt.Errorf("incorrect TG_ADMIN_ID value: %v", tgAdminId)
	}

	jwtAccessTTL, err := time.ParseDuration(values["JWT_ACCESS_TTL"])
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL value %q: %w", values["JWT_ACCESS_TTL"], err)
	}

	jwtRefreshTTL, err := time.ParseDuration(values["JWT_REFRESH_TTL"])
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL value %q: %w", values["JWT_REFRESH_TTL"], err)
	}

	return &Config{
		DbHost:                    values["DB_HOST"],
		DbPort:                    values["DB_PORT"],
		DbUser:                    values["DB_USER"],
		DbPassword:                values["DB_PASSWORD"],
		DbName:                    values["DB_NAME"],
		DbSSLMode:                 values["DB_SSLMODE"],
		RedisHost:                 values["REDIS_HOST"],
		RedisPort:                 values["REDIS_PORT"],
		ServerHost:                values["SERVER_HOST"],
		ServerPort:                values["SERVER_PORT"],
		ImageStoragePath:          values["IMAGE_STORAGE_PATH"],
		NotifyServiceURL:          values["NOTIFY_SERVICE_URL"],
		ExternalURL:               values["EXTERNAL_URL"],
		DefaultLogsPath:           values["DEFAULT_LOGS_PATH"],
		CriticalLogsPath:          values["CRITICAL_LOGS_PATH"],
		LoggerDebug:               loggerDebug,
		JwtSecret:                 values["JWT_SECRET"],
		JwtAccessTTL:              jwtAccessTTL,
		JwtRefreshTTL:             jwtRefreshTTL,
		PaginationDefaultPage:     0,
		PaginationDefaultPageSize: 10,
		TgAdminId:                 tgAdminId,
	}, nil
}
