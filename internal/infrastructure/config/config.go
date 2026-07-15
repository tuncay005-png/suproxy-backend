package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Log         LogConfig      `mapstructure:"log"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Xray        XrayConfig     `mapstructure:"xray"`
	Metrics     MetricsConfig  `mapstructure:"metrics"`
}

type ServerConfig struct {
	Address         string `mapstructure:"address"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type JWTConfig struct {
	SecretKey          string `mapstructure:"secret_key"`
	AccessTokenExpiry  int    `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry int    `mapstructure:"refresh_token_expiry"`
	Issuer             string `mapstructure:"issuer"`
}

type XrayConfig struct {
	UseMock    bool   `mapstructure:"use_mock"`
	BinaryPath string `mapstructure:"binary_path"`
	ConfigDir  string `mapstructure:"config_dir"`
	LogDir     string `mapstructure:"log_dir"`
	BackupDir  string `mapstructure:"backup_dir"`
	InstallDir string `mapstructure:"install_dir"`
}

type MetricsConfig struct {
	Enabled            bool `mapstructure:"enabled"`
	CollectionInterval int  `mapstructure:"collection_interval"` // in seconds
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("SUPROXY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Default values
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.shutdown_timeout", 10)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "suproxy")
	viper.SetDefault("database.password", "suproxy")
	viper.SetDefault("database.dbname", "suproxy")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", 5)
	viper.SetDefault("database.conn_max_idle_time", 10)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("jwt.secret_key", "change-me-in-production")
	viper.SetDefault("jwt.access_token_expiry", 15)
	viper.SetDefault("jwt.refresh_token_expiry", 168)
	viper.SetDefault("jwt.issuer", "suproxy-backend")
	viper.SetDefault("xray.use_mock", true)
	viper.SetDefault("xray.binary_path", "/usr/local/bin/xray")
	viper.SetDefault("xray.config_dir", "/etc/xray")
	viper.SetDefault("xray.log_dir", "/var/log/xray")
	viper.SetDefault("xray.backup_dir", "/var/backups/xray")
	viper.SetDefault("xray.install_dir", "/opt/xray")
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.collection_interval", 30)
}
