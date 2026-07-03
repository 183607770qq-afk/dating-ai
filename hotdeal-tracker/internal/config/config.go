package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Crawler  CrawlerConfig  `mapstructure:"crawler"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type CrawlerConfig struct {
	UserAgent        string            `mapstructure:"user_agent"`
	Timeout          int               `mapstructure:"timeout"`
	MaxDepth         int               `mapstructure:"max_depth"`
	Delay            int               `mapstructure:"delay"`
	Concurrent       int               `mapstructure:"concurrent"`
	ProxyPool        []string          `mapstructure:"proxy_pool"`
	RetryTimes       int               `mapstructure:"retry_times"`
	EnableHeadless   bool              `mapstructure:"enable_headless"`
	Platforms        map[string]bool   `mapstructure:"platforms"`
	HotKeywords      []string          `mapstructure:"hot_keywords"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/hotdeal-tracker/")
	viper.AddConfigPath("$HOME/.hotdeal-tracker")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 60)
	viper.SetDefault("server.write_timeout", 60)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("crawler.user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	viper.SetDefault("crawler.timeout", 30)
	viper.SetDefault("crawler.max_depth", 3)
	viper.SetDefault("crawler.delay", 1000)
	viper.SetDefault("crawler.concurrent", 5)
	viper.SetDefault("crawler.retry_times", 3)
	viper.SetDefault("crawler.enable_headless", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
