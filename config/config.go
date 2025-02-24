package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config 存储数据库和 JWT 的配置信息
type Config struct {
	DBDriver   string `env:"DB_DRIVER" envDefault:"sqlite3"` // 数据库驱动，可选 "mysql" 或 "sqlite3"
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"3306"`
	DBUser     string `env:"DB_USER" envDefault:"root"`
	DBPassword string `env:"DB_PASSWORD" envDefault:""`
	DBName     string `env:"DB_NAME" envDefault:"dlm"`
	SQLitePath string `env:"SQLITE_PATH" envDefault:"sqlite.db"` // sqlite数据库文件路径
	JWTSecret  string `env:"JWT_SECRET" envDefault:"secret"`
}

// DSN 返回数据库连接字符串，根据驱动不同返回不同的DSN
func (cfg *Config) DSN() string {
	if cfg.DBDriver == "sqlite3" {
		return cfg.SQLitePath
	} else if cfg.DBDriver == "mysql" {
		// MySQL: username:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	} else {
		panic("unknown db driver")
	}
}

// LoadConfig 解析环境变量并返回配置信息
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
