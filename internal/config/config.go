package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OTP      OTPConfig
	RateLimit RateLimitConfig
	SMTP     SMTPConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	Debug   bool
}

type MongoDBConfig struct {
	URI         string
	Database    string
	MaxPoolSize int
	MinPoolSize int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret                   string
	AccessTokenExpireMinutes int
	RefreshTokenExpireHours  int
}

type OTPConfig struct {
	Length        int
	ExpireMinutes int
}

type RateLimitConfig struct {
	RequestsPerMinute      int
	LoginRequestsPerMinute int
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v. Using environment variables.", err)
	}

	// Set defaults
	setDefaults()

	config := &Config{
		App: AppConfig{
			Name:  viper.GetString("APP_NAME"),
			Env:   viper.GetString("APP_ENV"),
			Port:  viper.GetString("APP_PORT"),
			Debug: viper.GetBool("APP_DEBUG"),
		},
		MongoDB: MongoDBConfig{
			URI:         viper.GetString("MONGODB_URI"),
			Database:    viper.GetString("MONGODB_DATABASE"),
			MaxPoolSize: viper.GetInt("MONGODB_MAX_POOL_SIZE"),
			MinPoolSize: viper.GetInt("MONGODB_MIN_POOL_SIZE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
		},
		JWT: JWTConfig{
			Secret:                   viper.GetString("JWT_SECRET"),
			AccessTokenExpireMinutes: viper.GetInt("JWT_ACCESS_TOKEN_EXPIRE_MINUTES"),
			RefreshTokenExpireHours:  viper.GetInt("JWT_REFRESH_TOKEN_EXPIRE_HOURS"),
		},
		OTP: OTPConfig{
			Length:        viper.GetInt("OTP_LENGTH"),
			ExpireMinutes: viper.GetInt("OTP_EXPIRE_MINUTES"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute:      viper.GetInt("RATE_LIMIT_REQUESTS_PER_MINUTE"),
			LoginRequestsPerMinute: viper.GetInt("RATE_LIMIT_LOGIN_REQUESTS_PER_MINUTE"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods: viper.GetStringSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders: viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		},
	}

	return config, nil
}

func setDefaults() {
	viper.SetDefault("APP_NAME", "Pakistani ERP System")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_DEBUG", true)

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGODB_DATABASE", "pakistani_erp")
	viper.SetDefault("MONGODB_MAX_POOL_SIZE", 100)
	viper.SetDefault("MONGODB_MIN_POOL_SIZE", 10)

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_POOL_SIZE", 10)

	viper.SetDefault("JWT_SECRET", "change-this-secret-key")
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", 15)
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRE_HOURS", 168)

	viper.SetDefault("OTP_LENGTH", 6)
	viper.SetDefault("OTP_EXPIRE_MINUTES", 5)

	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 10)
	viper.SetDefault("RATE_LIMIT_LOGIN_REQUESTS_PER_MINUTE", 5)
}

func (c *Config) GetAccessTokenDuration() time.Duration {
	return time.Duration(c.JWT.AccessTokenExpireMinutes) * time.Minute
}

func (c *Config) GetRefreshTokenDuration() time.Duration {
	return time.Duration(c.JWT.RefreshTokenExpireHours) * time.Hour
}

func (c *Config) GetOTPExpiration() time.Duration {
	return time.Duration(c.OTP.ExpireMinutes) * time.Minute
}
